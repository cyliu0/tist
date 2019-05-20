package matrix

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_PermuatateMatrix(t *testing.T) {
	matrix := make([][]interface{}, 3)
	for i := 0; i < 3; i++ {
		line := make([]interface{}, 3)
		for j := 0; j < 3; j++ {
			line[j] = fmt.Sprintf("%d", i*3+j)
		}
		matrix[i] = line
	}
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
			//existKey = fmt.Sprintf("%s%02d", existKey, item.Value)
			existKey = existKey + item.Value.(string)
		}
		if v, ok := existMap[existKey]; ok {
			t.Fatalf("Existed key: %v, count: %v,  current: %v", existKey, v, count)
		}
		logrus.Infof("iter.ID: %v, existKey: %v", iter.ID, existKey)
		existMap[existKey] = count
	}
}
