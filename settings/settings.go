package settings

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
