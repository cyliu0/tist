package sql

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/cyliu0/tist/matrix"
	"github.com/sirupsen/logrus"
)

type SQLMatrix struct {
	*matrix.PermutateMatrix
}

func NewSQLMatrix(sqlFilePrefix, sqlFilePostfix string, clientNumber int) (sqlMatrix *SQLMatrix, err error) {
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
		return
	}
	return
}

func getSQLFileNames(sqlFilePrefix, sqlFilePostfix string, clientNumber int) []string {
	fileNames := make([]string, clientNumber)
	for i := 0; i < clientNumber; i++ {
		fileNames[i] = fmt.Sprintf("%s-%d.%s", sqlFilePrefix, i, sqlFilePostfix)
	}
	return fileNames
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 1024))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}
