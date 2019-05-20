package worker

import (
	"github.com/cyliu0/tist/client"
	"github.com/cyliu0/tist/matrix"
	"github.com/cyliu0/tist/sql"
	"github.com/sirupsen/logrus"
)

type SQLMatrix struct {
	matrix.PermutateMatrix
}

func worker(tidbID int, clients []client.Client, iterChan <-chan matrix.Iterator) {
	for iterator := range iterChan {
		for sqlItem, err := iterator.NextItem(); err == nil; sqlItem, err = iterator.NextItem() {
			clientID := sqlItem.LineNum
			execStr := sqlItem.Value
			err := clients[clientID].Exec(tidbID, execStr.(string))
			if err != nil {
				logrus.Fatalf("Failed to execute sql: %v, err: %v", execStr, err)
			}
		}
	}
}

func RunWithWorkers(tidbNumber int, clients []client.Client, sqlMatrix *sql.SQLMatrix) {
	iterChans := make([]chan matrix.Iterator, tidbNumber)
	for tidbID := 0; tidbID < tidbNumber; tidbID++ {
		iterChan := make(chan matrix.Iterator, 1)
		iterChans[tidbID] = iterChan
		go worker(tidbID, clients, iterChans[tidbID])
	}
	for iterator, err := sqlMatrix.NextIterator(); err == nil; iterator, err = sqlMatrix.NextIterator() {
		iterChans[iterator.ID%tidbNumber] <- iterator
	}
}

func RunSync(tidbNumber int, clients []client.Client, sqlMatrix *sql.SQLMatrix) {
	for iterator, err := sqlMatrix.NextIterator(); err == nil; iterator, err = sqlMatrix.NextIterator() {
		tidbID := iterator.ID % tidbNumber
		for sqlItem, err := iterator.NextItem(); err == nil; sqlItem, err = iterator.NextItem() {
			clientID := sqlItem.LineNum
			execStr := sqlItem.Value
			err := clients[clientID].Exec(tidbID, execStr.(string))
			if err != nil {
				logrus.Fatalf("Failed to execute sql: %v, err: %v", execStr, err)
			}
		}
	}
}
