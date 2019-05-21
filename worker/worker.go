package worker

import (
	"github.com/cyliu0/tist/client"
	"github.com/cyliu0/tist/matrix"
	"github.com/cyliu0/tist/sql"
	"github.com/sirupsen/logrus"
)

func worker(tidbID int, clients []client.Client, iterChan <-chan matrix.Iterator, doneChan chan<- bool) {
	for iterator := range iterChan {
		for sqlItem, err := iterator.NextItem(); err == nil; sqlItem, err = iterator.NextItem() {
			clientID := sqlItem.LineNum
			execStr := sqlItem.Value
			err := clients[clientID].Exec(tidbID, execStr.(string))
			if err != nil {
				logrus.Fatalf("Failed to execute sql: %v, err: %v", execStr, err)
			}
		}
		logrus.Debugf("Iterator ID %d execution finished with TiDB ID: %d", iterator.ID, tidbID)
		doneChan <- true
	}
}

func RunWithWorkers(tidbNumber int, clients []client.Client, sqlMatrix *sql.SQLMatrix) {
	iterChans := make([]chan matrix.Iterator, tidbNumber)
	for tidbID := 0; tidbID < tidbNumber; tidbID++ {
		iterChan := make(chan matrix.Iterator, 1)
		iterChans[tidbID] = iterChan
		go worker(tidbID, clients, iterChan, sqlMatrix.DoneChan)
	}

	for iterator, err := sqlMatrix.NextIterator(); err == nil; iterator, err = sqlMatrix.NextIterator() {
		logrus.Debugf("Get iterator: %d", iterator.ID)
		workerID := iterator.ID % tidbNumber
		iterChans[workerID] <- iterator
		logrus.Debugf("worker ID: %d sending SQL with iterator ID: %d", workerID, iterator.ID)
		if (iterator.ID+1)%1000 == 0 {
			logrus.Infof("Working on No.%d permutation", iterator.ID+1)
		}
	}
}
