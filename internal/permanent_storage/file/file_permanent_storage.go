package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"os"
)

type PermanentStorage struct {
	FileFullPath string
	file         *os.File
	writer       *bufio.Writer
}

func (s *PermanentStorage) Import() ([]dto.Metrics, error) {
	if nil != s.file {
		s.Close()
		panic(fmt.Sprintf("File \"%s\" is open for write", s.FileFullPath))
	}

	file, scanner, err := s.openFileToRead()
	if errors.Is(err, os.ErrNotExist) {
		return []dto.Metrics{}, nil
	}
	if nil != err {
		return nil, err
	}

	defer file.Close()

	var metricList []dto.Metrics

	for scanner.Scan() {
		var metric dto.Metrics

		metricJson := scanner.Bytes()
		err = json.Unmarshal(metricJson, &metric)
		if nil != err {
			return nil, err
		}

		metricList = append(metricList, metric)
	}

	return metricList, nil
}

func (s *PermanentStorage) Export(metricList []dto.Metrics) error {
	err := s.openFileToWrite()
	if nil != err {
		return err
	}

	for _, metric := range metricList {
		metricJson, err := json.Marshal(metric)
		if nil != err {
			return err
		}

		err = s.writeMetricToFile(metricJson)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PermanentStorage) Close() {
	_ = s.file.Close()
}

func (s *PermanentStorage) writeMetricToFile(metricJson []byte) error {
	// Запись данных в буфер
	if _, err := s.writer.Write(metricJson); err != nil {
		return err
	}

	// Добавление переноса строки
	if err := s.writer.WriteByte('\n'); err != nil {
		return err
	}

	// Запись буфера в файл
	return s.writer.Flush()
}

func (s *PermanentStorage) openFileToRead() (*os.File, *bufio.Scanner, error) {
	file, err := os.OpenFile(s.FileFullPath, os.O_RDONLY, 0666)
	if nil != err {
		return file, &bufio.Scanner{}, err
	}

	return file, bufio.NewScanner(file), err
}

func (s *PermanentStorage) openFileToWrite() error {
	file, err := os.OpenFile(s.FileFullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if nil != err {
		return err
	}

	s.file = file
	s.writer = bufio.NewWriter(s.file)

	return nil
}

func NewInstance(fileFullPath string) interfaces.PermanentStorage {
	return &PermanentStorage{
		FileFullPath: fileFullPath,
	}
}
