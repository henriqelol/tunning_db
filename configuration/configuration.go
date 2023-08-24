package configuration

import (
	"errors"
	"strconv"
)

type Configuration struct {
	MemoryGB  string `json:"memory_gb"`
	CPUs      int    `json:"cpus"`
	DBType    string `json:"db_type"`
	DBVersion string `json:"db_version"`
}

func ParseRAM(ramString string) (int, error) {
	ramGBString := ramString[:len(ramString)-3]
	ramGB, err := strconv.Atoi(ramGBString)
	if err != nil {
		return 0, errors.New("field 'memory_db' should be specified in GB, like 'X GB'. Example '1 GB'")
	}
	return ramGB, nil
}
