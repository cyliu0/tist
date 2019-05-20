package matrix

import (
	"errors"

	"github.com/fighterlyt/permutation"
	"github.com/sirupsen/logrus"
)

type PermutateMatrix struct {
	iterSeqLength       int
	idle                chan bool
	iterChan            chan Iterator
	lineOrderPermutator *permutation.Permutator
	Matrix              [][]interface{}
}

type MatrixLocation struct {
	LineNum int
	Index   int
}

type IterItem struct {
	LineNum int
	Index   int
	Value   interface{}
}

type IterSeq []MatrixLocation

func (iterSeq IterSeq) GetMoveIndex() int {
	for i := len(iterSeq) - 1; i >= 1; i-- {
		if iterSeq[i].LineNum != iterSeq[i-1].LineNum {
			return i
		}
	}
	return 0
}

type Iterator struct {
	itemChan chan IterItem
	ID       int
}

func (iter Iterator) Iterate(pm *PermutateMatrix, iterSeq IterSeq) {
	for _, v := range iterSeq {
		iterItem := IterItem{
			LineNum: v.LineNum,
			Index:   v.Index,
			Value:   pm.Matrix[v.LineNum][v.Index],
		}
		iter.itemChan <- iterItem
	}
	close(iter.itemChan)
}

func NewIterator(iteratorID int, pm *PermutateMatrix, iterSeq IterSeq) Iterator {
	iter := Iterator{
		itemChan: make(chan IterItem, 1),
		ID:       iteratorID,
	}
	newIterSeq := make(IterSeq, len(iterSeq))
	copy(newIterSeq, iterSeq)
	// logrus.Infof("IteratorID: %v, iterSeq: %v", iteratorID, newIterSeq)
	go iter.Iterate(pm, newIterSeq)
	return iter
}

func (iter Iterator) NextItem() (*IterItem, error) {
	item, ok := <-iter.itemChan
	if !ok {
		return nil, errors.New("No next item")
	}
	return &item, nil
}

func NewPermutateMatrix(matrix [][]interface{}) (*PermutateMatrix, error) {
	pm := &PermutateMatrix{
		Matrix: matrix,
	}
	err := pm.newLineOrderPermutator()
	if err != nil {
		logrus.Errorf("Failed to initialize line order permutator, err: %v", err)
		return nil, err
	}
	pm.iterChan = make(chan Iterator, 1)
	go pm.iteratorGenerator()
	pm.idle = make(chan bool, 1)
	pm.idle <- true
	return pm, nil
}

func (pm *PermutateMatrix) newIterationSequence(lineOrder []int) IterSeq {
	iterSeq := IterSeq{}
	iterSeq = make([]MatrixLocation, 0)
	for _, lineNumber := range lineOrder {
		for index := 0; index < len(pm.Matrix[lineNumber]); index++ {
			iterSeq = append(iterSeq, MatrixLocation{
				LineNum: lineNumber,
				Index:   index,
			})
		}
	}
	pm.iterSeqLength = len(iterSeq)
	return iterSeq
}

func (pm *PermutateMatrix) newLineOrderPermutator() (err error) {
	lineNum := len(pm.Matrix)
	lineOrder := make([]int, lineNum)
	for i := 0; i < lineNum; i++ {
		lineOrder[i] = i
	}
	pm.lineOrderPermutator, err = permutation.NewPerm(lineOrder, nil)
	if err != nil {
		logrus.Errorf("Failed to initialize new permutator, err: %v", err)
	}
	return
}

func (pm *PermutateMatrix) nextLineOrder() ([]int, error) {
	lineOrder, err := pm.lineOrderPermutator.Next()
	if err != nil {
		if err.Error() == "all Permutations generated" {
			return nil, nil
		}
		logrus.Errorf("Failed to permutate next, err: %v", err)
		return nil, err
	}
	if lineOrder != nil {
		return lineOrder.([]int), nil
	}
	return nil, errors.New("Empty line order")
}

func (pm *PermutateMatrix) iteratorGenerator() {
	iteratorID := 0
	for lineOrder, err := pm.nextLineOrder(); lineOrder != nil && err == nil; lineOrder, err = pm.nextLineOrder() {
		iterSeq := pm.newIterationSequence(lineOrder)
		iter := NewIterator(iteratorID, pm, iterSeq)
		pm.iterChan <- iter
		iteratorID++
	L:
		for count := 1; count < pm.iterSeqLength; count++ {
			moveIndex := iterSeq.GetMoveIndex()
			if iterSeq[moveIndex].LineNum == iterSeq[0].LineNum {
				break
			}
			for i := moveIndex; i >= 2; i-- {
				if iterSeq[i].LineNum != iterSeq[i-1].LineNum {
					iterSeq[i-1], iterSeq[i] = iterSeq[i], iterSeq[i-1]
					iter = NewIterator(iteratorID, pm, iterSeq)
					iteratorID++
					pm.iterChan <- iter
					if (iterSeq[i-1].LineNum == iterSeq[0].LineNum) && (iterSeq[i-1].Index == i-1) {
						break L
					}
				}
			}
		}
	}
	close(pm.iterChan)
}

func (pm *PermutateMatrix) NextIterator() (iter Iterator, err error) {
	<-pm.idle
	iter, ok := <-pm.iterChan
	if !ok {
		err = errors.New("No next iteration sequence")
	}
	pm.idle <- true
	return
}
