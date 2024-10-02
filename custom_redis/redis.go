package custom_redis

import (
	"CustomRedis/common"
	"bufio"
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
	Rds.Restore()
	go Rds.BackgroundCleanupService()
}

func (r *Redis) Set(key, value string, ttl time.Duration, log bool) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[key] = RedisValue{value, time.Now().Add(ttl)}
	if log {
		err := r.Log(SET_OPERATION, key, value, ttl)
		if err != nil {
			fmt.Println("redis log error:", err)
			panic(err)
		}
	}
	return nil
}

func (r *Redis) Get(key string, log bool) (string, error) {
	r.lock.RLock()
	redisValue, ok := r.data[key]
	r.lock.RUnlock()

	if !ok {
		return "", KeyNotFoundError()
	}

	if time.Now().After(redisValue.ttl) {
		err := r.Delete(key, log)
		if err != nil {
			fmt.Println("redis ttl delete error:", err)
			panic(err)
		}
		return "", KeyNotFoundError()
	}

	return redisValue.value, nil
}

func (r *Redis) Delete(key string, log bool) error {
	r.lock.RLock()
	_, ok := r.data[key]
	r.lock.RUnlock()

	if !ok {
		return KeyNotFoundError()
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.data, key)
	if log {
		err := r.Log(DELETE_OPERATION, key, "", 0)
		if err != nil {
			// panic since restore will fail lead to data in consistency
			panic(err)
		}
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
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return errors.New(fmt.Sprintf("Error marshaling JSON : %s", err))
	}

	_, err = r.logfile.Write(append(jsonData, '\n'))
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Restore() error {
	scanner := bufio.NewScanner(r.logfile)
	fmt.Println("RESTORE LOG")
	for scanner.Scan() {
		var entry LogEntry
		line := scanner.Text()

		// Decode the JSON data
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return errors.New(fmt.Sprintf("Error unmarshaling JSON : %s", err))
		}

		// Print the parsed data
		if entry.Operation == SET_OPERATION {
			// remaining time for key
			remainingTtlDuration := time.Now().Sub(entry.CreatedAt.Add(entry.Ttl))
			if remainingTtlDuration > time.Second {
				fmt.Printf("SET Key : %s, Value : %s\n", entry.Key, entry.Value)
				err := r.Set(entry.Key, entry.Value, entry.Ttl, false)
				if err != nil {
					return err
				}
			}
		} else if entry.Operation == DELETE_OPERATION {
			err := r.Delete(entry.Key, false)
			if err != nil {
				return err
			}
			fmt.Printf("DELETE : %s\n", entry.Key)
		}
	}
	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return errors.New(fmt.Sprintf("Error reading redis log file for restore: %s", err))
	}
	return nil
}

func (r *Redis) BackgroundCleanupService() {
	for {
		for key := range r.data {
			if time.Now().After(r.data[key].ttl) {
				fmt.Printf("Deleting key : %s\n", key)
				err := r.Delete(key, true)
				if err != nil {
					fmt.Println("redis delete error:", err)
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
