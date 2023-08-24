package response

import (
	"fmt"
	"tuning_db/tuning"
)

func GenerateResponseMessage(tuningConfig tuning.TuningConfig) string {
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
