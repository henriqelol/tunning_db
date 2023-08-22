package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
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

// Flavor represents the input data for the tuning calculation
type Flavor struct {
	MemoryGB  int    `json:"memory_gb"`
	VCPUs     int    `json:"vcpus"`
	DBType    string `json:"db_type"`
	DBVersion string `json:"db_version"`
}

// TuningConfig represents the calculated tuning parameters
type TuningConfig struct {
	MaxConnections        int    `json:"max_connections"`
	InnoDBBufferPoolSize  string `json:"innodb_buffer_pool_size"`
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

func calculateMySQLTuning(flavor Flavor) TuningConfig {
	ramGB := float64(flavor.MemoryGB)

	maxConnections := int(math.Floor(ramGB/TNG_RAM_HALF_AVAILABLE) * TNG_MAX_CONNECTIONS_PER_VCPU)
	innodbBufferPoolSize := fmt.Sprintf("%.0fG", ramGB*TNG_INNODB_BUFFER_POOL_SIZE_RATIO)

	readBufferSize := int(math.Floor(float64(ramGB) * 0.1))
	if readBufferSize < 1 {
		readBufferSize = 1
	}

	tuningConfig := TuningConfig{
		MaxConnections:        maxConnections,
		InnoDBBufferPoolSize:  innodbBufferPoolSize,
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

	return tuningConfig
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var flavor Flavor
	err := json.NewDecoder(r.Body).Decode(&flavor)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tuningConfig := calculateMySQLTuning(flavor)
	json.NewEncoder(w).Encode(tuningConfig)
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/tune", handleRequest)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
