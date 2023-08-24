package tuning

import (
	"fmt"
	"math"
	"tuning_db/settings"
	"tuning_db/configuration"
)

type TuningConfig struct {
	MaxConnections        int    `json:"max_connections"`
	InnodbBufferPoolSize  string `json:"innodb_buffer_pool_size"`
	InnodbDedicatedServer string `json:"innodb_dedicated_server"`
	InnodbChangeBuffering string `json:"innodb_change_buffering"`
	SlowQueryLog          string `json:"slow_query_log"`
	SlowQueryLogFile      string `json:"slow_query_log_file"`
	LongQueryTime         int    `json:"long_query_time"`
	PerformanceSchema     int    `json:"performance_schema"`
	MaxAllowedPacket      string `json:"max_allowed_packet"`
	TableOpenCache        int    `json:"table_open_cache"`
	ThreadStack           string `json:"thread_stack"`
	ThreadCacheSize       int    `json:"thread_cache_size"`
	JoinBufferSize        string `json:"join_buffer_size"`
	SortBufferSize        string `json:"sort_buffer_size"`
	TmpTableSize          string `json:"tmp_table_size"`
	MaxHeapTableSize      string `json:"max_heap_table_size"`
	InnodbFlushMethod     string `json:"innodb_flush_method"`
	InnodbFilePerTable    int    `json:"innodb_file_per_table"`
	InnodbOpenFiles       int    `json:"innodb_open_files"`
	InnodbIoCapacity      int    `json:"innodb_io_capacity"`
	InnodbIoCapacityMax   int    `json:"innodb_io_capacity_max"`
	ReadBufferSize        int    `json:"read_buffer_size"`
}

func CalculateMySQLTuning(config configuration.Configuration) (TuningConfig, error) {
	ramGB, err := configuration.ParseRAM(config.MemoryGB)
	if err != nil {
		return TuningConfig{}, err
	}

	maxConnections := int(float64(ramGB) / settings.TNG_RAM_HALF_AVAILABLE * settings.TNG_MAX_CONNECTIONS_PER_VCPU)
	innodbBufferPoolSize := fmt.Sprintf("%dG", int(math.Floor(float64(ramGB)*settings.TNG_INNODB_BUFFER_POOL_SIZE_RATIO)))

	readBufferSize := int(math.Floor(float64(ramGB) * 0.1))
	if readBufferSize < 1 {
		readBufferSize = 1
	}

	tuningConfig := TuningConfig{
		MaxConnections:        maxConnections,
		InnodbBufferPoolSize:  innodbBufferPoolSize,
		InnodbDedicatedServer: settings.TNG_INNODB_DEDICATED_SERVER,
		InnodbChangeBuffering: settings.TNG_INNODB_CHANGE_BUFFERING,
		SlowQueryLog:          settings.TNG_SLOW_QUERY_LOG,
		SlowQueryLogFile:      settings.TNG_SLOW_QUERY_LOG_FILE,
		LongQueryTime:         settings.TNG_LONG_QUERY_TIME,
		PerformanceSchema:     settings.TNG_PERFORMANCE_SCHEMA,
		MaxAllowedPacket:      settings.TNG_MAX_ALLOWED_PACKET,
		TableOpenCache:        settings.TNG_TABLE_OPEN_CACHE,
		ThreadStack:           settings.TNG_THREAD_STACK,
		ThreadCacheSize:       settings.TNG_THREAD_CACHE_SIZE,
		JoinBufferSize:        settings.TNG_JOIN_BUFFER_SIZE,
		SortBufferSize:        settings.TNG_SORT_BUFFER_SIZE,
		TmpTableSize:          settings.TNG_TMP_TABLE_SIZE,
		MaxHeapTableSize:      settings.TNG_MAX_HEAP_TABLE_SIZE,
		InnodbFlushMethod:     settings.TNG_INNODB_FLUSH_METHOD,
		InnodbFilePerTable:    settings.TNG_INNODB_FILE_PER_TABLE,
		InnodbOpenFiles:       settings.TNG_INNODB_OPEN_FILES,
		InnodbIoCapacity:      settings.TNG_INNODB_IO_CAPACITY,
		InnodbIoCapacityMax:   settings.TNG_INNODB_IO_CAPACITY_MAX,
		ReadBufferSize:        readBufferSize,
	}

	return tuningConfig, nil
}
