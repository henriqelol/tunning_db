package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

// Settings for tuning database
const (
	TNG_INNODB_BUFFER_POOL_SIZE_RATIO = 0.7
	TNG_MAX_CONNECTIONS_PER_VCPU      = 50
	TNG_RAM_HALF_AVAILABLE            = 2
	TNG_INNODB_DEDICATED_SERVER       = "ON"
	TNG_INNODB_CHANGE_BUFFERING       = "none"
	TNG_SLOW_QUERY_LOG                = "ON"
	TNG_SLOW_QUERY_LOG_FILE           = "/var/log/mysql/slow_queries.log"
	TNG_LONG_QUERY_TIME               = 10
	TNG_PERFORMANCE_SCHEMA            = 1
	TNG_MAX_ALLOWED_PACKET            = "64M"
	TNG_TABLE_OPEN_CACHE              = 2000
	TNG_THREAD_STACK                  = "256K"
	TNG_THREAD_CACHE_SIZE             = -1
	TNG_JOIN_BUFFER_SIZE              = "2M"
	TNG_SORT_BUFFER_SIZE              = "2M"
	TNG_TMP_TABLE_SIZE                = "16M"
	TNG_MAX_HEAP_TABLE_SIZE           = "16M"
	TNG_INNODB_FLUSH_METHOD           = "O_DIRECT"
	TNG_INNODB_FILE_PER_TABLE         = 1
	TNG_INNODB_OPEN_FILES             = -1
	TNG_INNODB_IO_CAPACITY            = 200
	TNG_INNODB_IO_CAPACITY_MAX        = 2 * TNG_INNODB_IO_CAPACITY
)

// Configuration represents the input data for the tuning calculation
type Configuration struct {
	MemoryGB  string `json:"memory_gb"`
	CPUs      int    `json:"cpus"`
	DBType    string `json:"db_type"`
	DBVersion string `json:"db_version"`
}

// TuningConfig represents the calculated tuning parameters
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

func parseRAM(ramString string) (int, error) {
	ramGBString := ramString[:len(ramString)-3]
	ramGB, err := strconv.Atoi(ramGBString)
	if err != nil {
		return 0, errors.New("field 'memory_db' should be specified in GB, like 'X GB'. Example '1 GB'")
	}
	return ramGB, nil
}

func calculateMySQLTuning(configuration Configuration) (TuningConfig, error) {
	ramGB, err := parseRAM(configuration.MemoryGB)
	if err != nil {
		return TuningConfig{}, err
	}

	maxConnections := int(float64(ramGB) / TNG_RAM_HALF_AVAILABLE * TNG_MAX_CONNECTIONS_PER_VCPU)
	innodbBufferPoolSize := fmt.Sprintf("%dG", int(math.Floor(float64(ramGB)*TNG_INNODB_BUFFER_POOL_SIZE_RATIO)))

	readBufferSize := int(math.Floor(float64(ramGB) * 0.1))
	if readBufferSize < 1 {
		readBufferSize = 1
	}

	tuningConfig := TuningConfig{
		MaxConnections:        maxConnections,
		InnodbBufferPoolSize:  innodbBufferPoolSize,
		InnodbDedicatedServer: TNG_INNODB_DEDICATED_SERVER,
		InnodbChangeBuffering: TNG_INNODB_CHANGE_BUFFERING,
		SlowQueryLog:          TNG_SLOW_QUERY_LOG,
		SlowQueryLogFile:      TNG_SLOW_QUERY_LOG_FILE,
		LongQueryTime:         TNG_LONG_QUERY_TIME,
		PerformanceSchema:     TNG_PERFORMANCE_SCHEMA,
		MaxAllowedPacket:      TNG_MAX_ALLOWED_PACKET,
		TableOpenCache:        TNG_TABLE_OPEN_CACHE,
		ThreadStack:           TNG_THREAD_STACK,
		ThreadCacheSize:       TNG_THREAD_CACHE_SIZE,
		JoinBufferSize:        TNG_JOIN_BUFFER_SIZE,
		SortBufferSize:        TNG_SORT_BUFFER_SIZE,
		TmpTableSize:          TNG_TMP_TABLE_SIZE,
		MaxHeapTableSize:      TNG_MAX_HEAP_TABLE_SIZE,
		InnodbFlushMethod:     TNG_INNODB_FLUSH_METHOD,
		InnodbFilePerTable:    TNG_INNODB_FILE_PER_TABLE,
		InnodbOpenFiles:       TNG_INNODB_OPEN_FILES,
		InnodbIoCapacity:      TNG_INNODB_IO_CAPACITY,
		InnodbIoCapacityMax:   TNG_INNODB_IO_CAPACITY_MAX,
		ReadBufferSize:        readBufferSize,
	}

	return tuningConfig, nil
}

func GenerateResponseMessage(tuningConfig TuningConfig) string {
	return fmt.Sprintf(`
	[mysqld]

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
	# For example: /etc/mysql/conf.d/tunning_dbs.cnf`, tuningConfig.MaxConnections, tuningConfig.InnodbBufferPoolSize,
		tuningConfig.InnodbDedicatedServer, tuningConfig.InnodbChangeBuffering, tuningConfig.SlowQueryLog,
		tuningConfig.LongQueryTime, tuningConfig.PerformanceSchema, tuningConfig.MaxAllowedPacket,
		tuningConfig.TableOpenCache, tuningConfig.ThreadStack, tuningConfig.ThreadCacheSize, tuningConfig.JoinBufferSize,
		tuningConfig.SortBufferSize, tuningConfig.TmpTableSize, tuningConfig.MaxHeapTableSize,
		tuningConfig.InnodbFlushMethod, tuningConfig.InnodbFilePerTable, tuningConfig.InnodbOpenFiles,
		tuningConfig.InnodbIoCapacity, tuningConfig.InnodbIoCapacityMax, tuningConfig.ReadBufferSize)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var configuration Configuration
	err := json.NewDecoder(r.Body).Decode(&configuration)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tuningConfig, err := calculateMySQLTuning(configuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseMessage := GenerateResponseMessage(tuningConfig)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, responseMessage)
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/tuning", handleRequest)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
