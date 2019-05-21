package sql

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cyliu0/tist/matrix"
	"github.com/sirupsen/logrus"
)

type SQLMatrix struct {
	*matrix.PermutateMatrix
	SQLNum int
}

func NewSQLMatrix(sqlFilePrefix, sqlFileSuffix string, clientNumber int) (sqlMatrix *SQLMatrix, err error) {
	sqlMatrix = &SQLMatrix{}
	mtx := make([][]interface{}, clientNumber)
	fileNames := getSQLFileNames(sqlFilePrefix, sqlFileSuffix, clientNumber)
	for i, fileName := range fileNames {
		lines, err := readNonBlankLines(fileName)
		if err != nil {
			logrus.Fatalf("Failed to read lines from file: %v, err: %v", fileName, err)
		}
		linesLen := len(lines)
		if linesLen == 0 {
			logrus.Fatalf("SQL file: %s is empty", fileName)
		}
		sqlMatrix.SQLNum = sqlMatrix.SQLNum + linesLen
		sqlStrSlice := make([]interface{}, linesLen)
		for j, sqlStr := range lines {
			sqlStrSlice[j] = sqlStr
		}
		mtx[i] = sqlStrSlice
	}
	sqlMatrix.PermutateMatrix, err = matrix.NewPermutateMatrix(mtx)
	if err != nil {
		logrus.Errorf("Failed to new permutate matrix, err: %v", err)
	}
	return
}

func (sqlMatrix *SQLMatrix) Brief() {
	for clientID, sqls := range sqlMatrix.Matrix {
		logrus.Infof("ClientID: %d, SQL Num: %d", clientID, len(sqls))
	}
	logrus.Infof("Total SQL Num: %d", sqlMatrix.SQLNum)
}

func getSQLFileNames(sqlFilePrefix, sqlFilePostfix string, clientNumber int) []string {
	fileNames := make([]string, clientNumber)
	for i := 0; i < clientNumber; i++ {
		fileNames[i] = fmt.Sprintf("%s-%d.%s", sqlFilePrefix, i, sqlFilePostfix)
	}
	return fileNames
}

func readNonBlankLines(fileName string) ([]string, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Errorf("open file error: %v", err)
		return nil, err
	}
	lines := make([]string, 0)
	for _, line := range strings.Split(string(content), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
