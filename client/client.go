package client

import (
	"github.com/sirupsen/logrus"
)

// Client is a simulated client for testing TiDB
type Client struct {
	ClientID int
	TiDBs    []*TiDB
}

func (c Client) Exec(tidbID int, execStr string) (err error) {
	_, err = c.TiDBs[tidbID].db.Exec(execStr)
	if err != nil {
		logrus.Errorf("Failed to exec sql: %s, TiDB: %v", execStr, c.TiDBs[tidbID])
	}
	//logrus.Debugf("Exec SQL: %s on TiDB: %d", execStr, tidbID)
	return
}

func InitClients(clientNumber int, sqlFilePrefix, tidbJsonFile string) (clients []Client, tidbNumber int, err error) {
	tidbs, err := GetTiDBs(tidbJsonFile)
	if err != nil {
		logrus.Errorf("Failed to get TiDB executor, err: %v", err)
		return
	}
	tidbNumber = len(tidbs)
	for i := 0; i < clientNumber; i++ {
		client := Client{
			ClientID: i,
			TiDBs:    tidbs,
		}
		for tidbID, tidb := range client.TiDBs {
			if err = tidb.Connect(); err != nil {
				logrus.Errorf("Failed to connect to TiDB %v, err: %v", tidbID, err)
				return
			}
		}
		clients = append(clients, client)
	}
	return
}
