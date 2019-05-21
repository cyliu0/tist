package matrix

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_PermuatateMatrix(t *testing.T) {
	matrix := make([][]interface{}, 3)
	matrix[0] = []interface{}{"1", "2"} //, "7"}
	matrix[1] = []interface{}{"3"}      //"4", "8"}
	matrix[2] = []interface{}{"5"}      //"6", "9"}
	pm, err := NewPermutateMatrix(matrix)
	if err != nil {
		t.Errorf("Failed to new permutate matrix, err: %v", err)
	}
	count := 0
	existMap := make(map[string]int)
	for iter, err := pm.NextIterator(); err == nil; iter, err = pm.NextIterator() {
		count++
		existKey := ""
		for item, err := iter.NextItem(); err == nil; item, err = iter.NextItem() {
			existKey = existKey + item.Value.(string)
		}
		if v, ok := existMap[existKey]; ok {
			t.Fatalf("Existed key: %v, count: %v,  current: %v", existKey, v, count)
		}
		logrus.Infof("iter.ID: %v, existKey: %v", iter.ID, existKey)
		existMap[existKey] = count
		pm.DoneChan <- true
	}
}
