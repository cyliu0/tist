package cmd

import (
	"time"

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
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
		clients, tidbNumber, err := client.InitClients(clientNumber, sqlFilePrefix, tidbJsonFile)
		if err != nil {
			logrus.Fatalf("Failed to initialize clients, err: %v", err)
		}
		sqlMatrix, err := sql.NewSQLMatrix(sqlFilePrefix, sqlFileSuffix, clientNumber)
		if err != nil {
			logrus.Fatalf("Failed to initialize SQL matrix, err: %v", err)
		}
		sqlMatrix.Brief()
		start := time.Now()
		worker.RunWithWorkers(tidbNumber, clients, sqlMatrix)
		end := time.Now()
		duration := end.Sub(start)
		logrus.Infof("SQL Num: %v, Permutation Num: %d, Time Duration: %s", sqlMatrix.SQLNum, sqlMatrix.IterTotal, duration)
	},
}

func Execute() {
	if err := tistCmd.Execute(); err != nil {
		logrus.Fatalf("Command err: %v", err)
	}
}

var clientNumber int
var sqlFilePrefix string
var sqlFileSuffix string
var tidbJsonFile string
var verbose bool

func init() {
	tistCmd.Flags().StringVarP(&tidbJsonFile, "tidb-config", "t", "./config/tidb-clusters.json", "TiDB clusters JSON file")
	tistCmd.Flags().StringVarP(&sqlFilePrefix, "sql-file-prefix", "p", "./config/sql", "prefix for SQL files")
	tistCmd.Flags().StringVarP(&sqlFileSuffix, "sql-file-suffix", "s", "sql", "suffix for SQL files")
	tistCmd.Flags().IntVarP(&clientNumber, "client-number", "c", 3, "number of client")
	tistCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
