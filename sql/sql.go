package sql

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cyliu0/tist/matrix"
	"github.com/sirupsen/logrus"
)

type SQLMatrix struct {
	*matrix.PermutateMatrix
}

func NewSQLMatrix(sqlFilePrefix, sqlFilePostfix string, clientNumber int) (sqlMatrix *SQLMatrix, err error) {
	sqlMatrix = &SQLMatrix{}
	mtx := make([][]interface{}, clientNumber)
	fileNames := getSQLFileNames(sqlFilePrefix, sqlFilePostfix, clientNumber)
	for i, fileName := range fileNames {
		lines, err := readLines(fileName)
		if err != nil {
			logrus.Fatalf("Failed to read lines from file: %v, err: %v", fileName, err)
		}
		sqlStrSlice := make([]interface{}, len(lines))
		for j, sqlStr := range lines {
			sqlStrSlice[j] = sqlStr
		}
		mtx[i] = sqlStrSlice
	}
	sqlMatrix.PermutateMatrix, err = matrix.NewPermutateMatrix(mtx)
	if err != nil {
		logrus.Errorf("Failed to new permutate matrix, err: %v", err)
	}
	return sqlMatrix, err
}

func getSQLFileNames(sqlFilePrefix, sqlFilePostfix string, clientNumber int) []string {
	fileNames := make([]string, clientNumber)
	for i := 0; i < clientNumber; i++ {
		fileNames[i] = fmt.Sprintf("%s-%d.%s", sqlFilePrefix, i, sqlFilePostfix)
	}
	return fileNames
}

func readLines(fileName string) (lines []string, err error) {
	lines = make([]string, 0)
	f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logrus.Errorf("open file error: %v", err)
		return
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.Errorf("read file line error: %v", err)
			return nil, err
		}
		line = strings.TrimRight(line, "\n")
		lines = append(lines, line)
	}
	return
}
