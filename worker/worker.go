package worker

import (
	"time"

	"github.com/cyliu0/tist/client"
	"github.com/cyliu0/tist/matrix"
	"github.com/cyliu0/tist/sql"
	"github.com/sirupsen/logrus"
)

func worker(tidbID int, clients []client.Client, iterChan <-chan matrix.Iterator, doneChan chan<- bool) {
	for iterator := range iterChan {
		start := time.Now()
		for sqlItem, err := iterator.NextItem(); err == nil; sqlItem, err = iterator.NextItem() {
			clientID := sqlItem.LineNum
			execStr := sqlItem.Value
			err := clients[clientID].Exec(tidbID, execStr.(string))
			if err != nil {
				logrus.Fatalf("Failed to execute sql: %v, err: %v", execStr, err)
			}
		}
		end := time.Now()
		duration := end.Sub(start)
		logrus.Infof("Iter %d finished on TiDB: %v in %s", iterator.ID, tidbID, duration.String())
		doneChan <- true
	}
}

func RunWithWorkers(tidbNumber int, clients []client.Client, sqlMatrix *sql.SQLMatrix) {
	start := time.Now()
	iterChans := make([]chan matrix.Iterator, tidbNumber)
	for tidbID := 0; tidbID < tidbNumber; tidbID++ {
		iterChan := make(chan matrix.Iterator, 1)
		iterChans[tidbID] = iterChan
		go worker(tidbID, clients, iterChan, sqlMatrix.DoneChan)
	}

	for iterator, err := sqlMatrix.NextIterator(); err == nil; iterator, err = sqlMatrix.NextIterator() {
		workerID := iterator.ID % tidbNumber
		iterChans[workerID] <- iterator
	}

	end := time.Now()
	duration := end.Sub(start)
	logrus.Infof("Permutation Num: %d, Time duration: %s", sqlMatrix.IterTotal, duration)
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
