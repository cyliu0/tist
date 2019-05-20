package cmd

import (
	"github.com/cyliu0/tist/client"
	"github.com/cyliu0/tist/sql"
	"github.com/cyliu0/tist/worker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tistCmd = &cobra.Command{
	Use:  "tist",
	Long: `A test tool for parallel running sql against multiple TiDB clusters.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients, tidbNumber, err := client.InitClients(clientNumber, sqlFilePrefix, tidbJsonFile)
		if err != nil {
			logrus.Fatalf("Failed to initialize clients, err: %v", err)
		}
		sqlMatrix, err := sql.NewSQLMatrix(sqlFilePrefix, sqlFilePostfix, clientNumber)
		if err != nil {
			logrus.Fatalf("Failed to initialize SQL matrix, err: %v", err)
		}
		worker.RunWithWorkers(tidbNumber, clients, sqlMatrix)
	},
}

func Execute() {
	if err := tistCmd.Execute(); err != nil {
		logrus.Fatalf("Command err: %v", err)
	}
}

var clientNumber int
var sqlFilePrefix string
var sqlFilePostfix string
var tidbJsonFile string

func init() {
	tistCmd.Flags().StringVar(&tidbJsonFile, "tidb-config", "./config/tidb-clusters.json", "TiDB clusters JSON file")
	tistCmd.Flags().StringVar(&sqlFilePrefix, "sql-file-prefix", "./config/sql", "Prefix for SQL files")
	tistCmd.Flags().StringVar(&sqlFilePostfix, "sql-file-postfix", "sql", "Postfix for SQL files")
	tistCmd.Flags().IntVar(&clientNumber, "client-number", 3, "Number of client")
}
