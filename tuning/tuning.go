package tuning

import (
	"fmt"
	"io/ioutil"
	"math"
	"tuning_db/configuration"
	"tuning_db/settings"
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

func calculateMySQLTuning(config configuration.Configuration) (TuningConfig, error) {
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

func generateArchiveTuning(tuningConfig TuningConfig) string {
	return fmt.Sprintf(`[mysqld]

max_connections=%d
innodb_buffer_pool_size=%s
innodb_dedicated_server=%s
innodb_change_buffering=%s
slow_query_log=%s
long_query_time=%d
performance_schema=%d
max_allowed_packet=%s
table_open_cache=%d
thread_stack=%s
thread_cache_size=%d
join_buffer_size=%s
sort_buffer_size=%s
tmp_table_size=%s
max_heap_table_size=%s
innodb_flush_method=%s
innodb_file_per_table=%d
innodb_open_files=%d
innodb_io_capacity=%d
innodb_io_capacity_max=%d
read_buffer_size=%d

# Save in a file with extension .cnf in the directory /etc/mysql/conf.d/
# For example: /etc/mysql/conf.d/tuning_db.cnf

# More details about variables in :https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html
	`, tuningConfig.MaxConnections, tuningConfig.InnodbBufferPoolSize,
		tuningConfig.InnodbDedicatedServer, tuningConfig.InnodbChangeBuffering, tuningConfig.SlowQueryLog,
		tuningConfig.LongQueryTime, tuningConfig.PerformanceSchema, tuningConfig.MaxAllowedPacket,
		tuningConfig.TableOpenCache, tuningConfig.ThreadStack, tuningConfig.ThreadCacheSize, tuningConfig.JoinBufferSize,
		tuningConfig.SortBufferSize, tuningConfig.TmpTableSize, tuningConfig.MaxHeapTableSize,
		tuningConfig.InnodbFlushMethod, tuningConfig.InnodbFilePerTable, tuningConfig.InnodbOpenFiles,
		tuningConfig.InnodbIoCapacity, tuningConfig.InnodbIoCapacityMax, tuningConfig.ReadBufferSize)
}

func ConfigureLocalTuningDatabase(config configuration.Configuration) error {
	tuningConfig, err := calculateMySQLTuning(config)
	if err != nil {
		return err
	}

	fileContent := generateArchiveTuning(tuningConfig)

	fmt.Print(fileContent)

	// Define o caminho completo do arquivo no diretório /etc/mysql/conf.d/
	filePath := "./tuning_db.cnf"

	// Escreve o conteúdo no arquivo
	err = ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		fmt.Print(err)
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Config file saved to: %s\n", filePath)
	return nil
}
