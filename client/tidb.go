package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// TiDB is a TiDB SQL executor.
type TiDB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password,omitempty"`
	Database string `json:"database"`
	db       *sql.DB
}

func (t *TiDB) GetConnectString() string {
	if t.Password == "" {
		return fmt.Sprintf("%s@tcp(%s:%d)/%s", t.User, t.Host, t.Port, t.Database)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", t.User, t.Password, t.Host, t.Port, t.Database)
}

func (t *TiDB) Connect() (err error) {
	source := t.GetConnectString()
	t.db, err = sql.Open("mysql", source)
	if err != nil {
		logrus.Errorf("Failed to connect to TiDB: %v, sqlx.Open err: %v\n", source, err)
	}
	return
}

func GetTiDBs(tidbJsonFile string) (TiDBs []*TiDB, err error) {
	content, err := ioutil.ReadFile(tidbJsonFile)
	if err != nil {
		logrus.Errorf("Failed to read from %v, ioutil.ReadFile err: %v", tidbJsonFile, err)
		return
	}
	err = json.Unmarshal(content, &TiDBs)
	if err != nil {
		logrus.Errorf("Failed to unmarshal json content: %s, json.Unmarshal err: %v", content, err)
	}
	return
}
