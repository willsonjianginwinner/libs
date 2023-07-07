package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-ini/ini"
)

var configMutex sync.RWMutex
var envConfig *ini.File

func init() {
	prefix := "./"

	err := Set(prefix, nil)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
		return
	}
}

func Set(prefix string, files []any) error {
	iniFile, err := loadConfig(prefix, files)

	if err != nil {
		return fmt.Errorf("fail to read file: %v", err)
	}

	configMutex.Lock()
	defer configMutex.Unlock()
	envConfig = iniFile
	return nil
}

func loadConfig(prefix string, files []any) (*ini.File, error) {

	for idx, val := range files {
		if config, ok := val.(string); ok {
			config = prefix + config
			files[idx] = config
		}
	}

	iniFile, err := ini.Load(prefix+"env.ini", files...)

	return iniFile, err
}

func Get(section, key, defaultValue string) string {
	configMutex.RLock()
	defer configMutex.RUnlock()

	var value = defaultValue
	if v := envConfig.Section(section).Key(key).String(); v != "" {
		value = v
	}

	return value
}

func GetBool(section, key string, defaultValue bool) bool {
	configMutex.RLock()
	defer configMutex.RUnlock()

	var value = defaultValue

	if v, err := envConfig.Section(section).Key(key).Bool(); nil == err {
		value = v
	}

	return value
}

func GetInt(section, key string, defaultValue int) int {
	configMutex.RLock()
	defer configMutex.RUnlock()

	var value = defaultValue

	if v, err := envConfig.Section(section).Key(key).Int(); nil == err {
		value = v
	}

	return value
}

func GetInt64(section, key string, defaultValue int64) int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()

	var value = defaultValue

	if v, err := envConfig.Section(section).Key(key).Int64(); nil == err {
		value = v
	}

	return value
}
