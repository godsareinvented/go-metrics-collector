package dictionary

const (
	AllocMetricName         = "Alloc"
	BuckHashSysMetricName   = "BuckHashSys"
	FreesMetricName         = "Frees"
	GCCPUFractionMetricName = "GCCPUFraction"
	GCSysMetricName         = "GCSys"
	HeapAllocMetricName     = "HeapAlloc"
	HeapIdleMetricName      = "HeapIdle"
	HeapInuseMetricName     = "HeapInuse"
	HeapObjectsMetricName   = "HeapObjects"
	HeapReleasedMetricName  = "HeapReleased"
	HeapSysMetricName       = "HeapSys"
	LastGCMetricName        = "LastGC"
	LookupsMetricName       = "Lookups"
	MCacheInuseMetricName   = "MCacheInuse"
	MCacheSysMetricName     = "MCacheSys"
	MSpanInuseMetricName    = "MSpanInuse"
	MSpanSysMetricName      = "MSpanSys"
	MallocsMetricName       = "Mallocs"
	NextGCMetricName        = "NextGC"
	NumForcedGCMetricName   = "NumForcedGC"
	NumGCMetricName         = "NumGC"
	OtherSysMetricName      = "OtherSys"
	PauseTotalNsMetricName  = "PauseTotalNs"
	StackInuseMetricName    = "StackInuse"
	StackSysMetricName      = "StackSys"
	SysMetricName           = "Sys"
	TotalAllocMetricName    = "TotalAlloc"
	PollCountMetricName     = "PollCount"
	RandomValueMetricName   = "RandomValue"

	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

var MetricNameList = [29]string{
	AllocMetricName,
	BuckHashSysMetricName,
	FreesMetricName,
	GCCPUFractionMetricName,
	GCSysMetricName,
	HeapAllocMetricName,
	HeapIdleMetricName,
	HeapInuseMetricName,
	HeapObjectsMetricName,
	HeapReleasedMetricName,
	HeapSysMetricName,
	LastGCMetricName,
	LookupsMetricName,
	MCacheInuseMetricName,
	MCacheSysMetricName,
	MSpanInuseMetricName,
	MSpanSysMetricName,
	MallocsMetricName,
	NextGCMetricName,
	NumForcedGCMetricName,
	NumGCMetricName,
	OtherSysMetricName,
	PauseTotalNsMetricName,
	StackInuseMetricName,
	StackSysMetricName,
	SysMetricName,
	TotalAllocMetricName,
	PollCountMetricName,
	RandomValueMetricName,
}
