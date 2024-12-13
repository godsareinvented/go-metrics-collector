package mem_storage

import (
	"context"
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"strconv"
	"sync"
)

type MemStorage struct {
	mu         sync.Mutex
	entityList [][]byte
	nameIndex  map[string]int
	idIndex    map[string]int
}

func (memStorage *MemStorage) GetAll(_ context.Context) ([]dto.Metrics, error) {
	var metricList []dto.Metrics

	for _, metricJson := range memStorage.entityList {
		metric, err := memStorage.getDecodedMetric(metricJson)
		if nil != err {
			return nil, err
		}

		metricList = append(metricList, metric)
	}

	return metricList, nil
}

// GetByID todo: ID - уникальный. Зачем поиск по типу?
func (memStorage *MemStorage) GetByID(_ context.Context, ID string, _ string) (dto.Metrics, bool, error) {
	index := memStorage.getMetricIndexByID(ID)
	if -1 == index {
		return dto.Metrics{}, false, nil
	}
	metric, err := memStorage.getDecodedMetric(memStorage.entityList[index])
	if nil != err {
		return dto.Metrics{}, false, err
	}
	return metric, true, nil
}

// GetByName todo: ...
func (memStorage *MemStorage) GetByName(_ context.Context, mName string, mType string) (dto.Metrics, bool, error) {
	index := memStorage.getMetricIndexByName(mName)
	if -1 == index {
		return dto.Metrics{}, false, nil
	}
	metric, err := memStorage.getDecodedMetric(memStorage.entityList[index])
	if nil != err {
		return dto.Metrics{}, false, err
	}
	return metric, true, nil
}

func (memStorage *MemStorage) Save(ctx context.Context, metric dto.Metrics) (string, error) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	if "" == metric.ID {
		ID, _ := memStorage.GetGeneratedID(ctx, metric)
		metric.ID = ID
	}

	metricJson, err := memStorage.getEncodedMetric(metric)
	if nil != err {
		return "", err
	}

	if index := memStorage.getMetricIndex(metric); -1 != index {
		memStorage.update(index, metric, metricJson)
		return strconv.Itoa(index), nil
	}

	index := memStorage.save(metric, metricJson)
	return strconv.Itoa(index), nil
}

func (memStorage *MemStorage) SaveBatch(ctx context.Context, metricBatch []dto.Metrics) error {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	for _, metric := range metricBatch {
		if "" == metric.ID {
			ID, _ := memStorage.GetGeneratedID(ctx, metric)
			metric.ID = ID
		}

		metricJson, err := memStorage.getEncodedMetric(metric)
		if nil != err {
			return err
		}
		if index := memStorage.getMetricIndex(metric); -1 != index {
			memStorage.update(index, metric, metricJson)
			strconv.Itoa(index)
			continue
		}

		memStorage.save(metric, metricJson)
	}

	return nil
}

// GetGeneratedID Метод предполагает, что все идентификаторы генерятся либо вне, либо внутри.
// Проверки на уникальность сгенерированного внутри идентификатора нет.
func (memStorage *MemStorage) GetGeneratedID(_ context.Context, metric dto.Metrics) (string, error) {
	if "" != metric.ID {
		return metric.ID, nil
	}

	if index := memStorage.getMetricIndex(metric); -1 != index {
		return strconv.Itoa(index), nil
	}

	return strconv.Itoa(len(memStorage.entityList)), nil
}

func (memStorage *MemStorage) save(metric dto.Metrics, metricJson []byte) int {
	memStorage.entityList = append(memStorage.entityList, metricJson)

	index := len(memStorage.entityList) - 1
	memStorage.idIndex[metric.ID] = index
	memStorage.nameIndex[metric.MName] = index

	return index
}

func (memStorage *MemStorage) update(index int, metric dto.Metrics, metricJson []byte) {
	memStorage.entityList[index] = metricJson

	oldMetricName := memStorage.getFoundName(index)
	if oldMetricName != metric.MName {
		delete(memStorage.nameIndex, oldMetricName)
		memStorage.nameIndex[metric.MName] = index
	}
}

func (memStorage *MemStorage) getFoundName(index int) string {
	for metricName, i := range memStorage.nameIndex {
		if index == i {
			return metricName
		}
	}
	return ""
}

func (memStorage *MemStorage) getMetricIndex(metric dto.Metrics) int {
	if index := memStorage.getMetricIndexByID(metric.ID); -1 != index {
		return index
	}
	if index := memStorage.getMetricIndexByName(metric.MName); -1 != index {
		return index
	}
	return -1
}

func (memStorage *MemStorage) getMetricIndexByID(ID string) int {
	if "" == ID {
		return -1
	}
	if index, ok := memStorage.idIndex[ID]; ok {
		return index
	}
	return -1
}

func (memStorage *MemStorage) getMetricIndexByName(metricName string) int {
	if "" == metricName {
		return -1
	}
	if index, ok := memStorage.nameIndex[metricName]; ok {
		return index
	}
	return -1
}

func (memStorage *MemStorage) getEncodedMetric(metric dto.Metrics) ([]byte, error) {
	metricJson, err := json.Marshal(metric)
	if nil != err {
		return metricJson, err
	}
	return metricJson, nil
}

func (memStorage *MemStorage) getDecodedMetric(metricJson []byte) (dto.Metrics, error) {
	metric := dto.Metrics{}
	err := json.Unmarshal(metricJson, &metric)
	if nil != err {
		return metric, err
	}
	return metric, nil
}

func NewInstance(idIndex map[string]int, nameIndex map[string]int) interfaces.StorageInterface {
	return &MemStorage{
		idIndex:   idIndex,
		nameIndex: nameIndex,
	}
}
