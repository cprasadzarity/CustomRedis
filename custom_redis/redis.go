package custom_redis

import (
	"CustomRedis/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"sync"
	"time"
)

var Rds *Redis

const (
	SET_OPERATION    = "set"
	DELETE_OPERATION = "del"
)

func Init() {
	godotenv.Load()
	redisLogFile := os.Getenv("LOG_FILE")
	file, err := common.OpenFile(redisLogFile)
	if err != nil {
		fmt.Println("open redis log file error:", err)
		panic(err)
	}
	Rds = &Redis{make(map[string]RedisValue), sync.RWMutex{}, file}
}

func (r *Redis) Set(key, value string, ttl time.Duration) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[key] = RedisValue{value, time.Now().Add(ttl)}
	err := r.Log(SET_OPERATION, key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key string) (string, error) {
	r.lock.RLock()
	redisValue, ok := r.data[key]
	r.lock.RUnlock()

	if !ok {
		return "", KeyNotFoundError()
	}

	if time.Now().After(redisValue.ttl) {
		r.Delete(key)
		return "", KeyNotFoundError()
	}

	return redisValue.value, nil
}

func (r *Redis) Delete(key string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.data, key)
	err := r.Log(DELETE_OPERATION, key, "", 0)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Log(operation, key, value string, ttl time.Duration) error {
	logEntry := LogEntry{
		Operation: operation,
		Key:       key,
		Value:     value,
		Ttl:       ttl,
		CreatedAt: time.Now(),
	}
	jsonData, err := json.MarshalIndent(logEntry, "", "  ") // Indented for better readability
	if err != nil {
		return errors.New(fmt.Sprintf("Error marshaling JSON : %s", err))
	}

	_, err = r.logfile.Write(append(jsonData, '\n'))
	if err != nil {
		return err
	}

	return nil
}
