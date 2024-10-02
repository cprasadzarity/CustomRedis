package custom_redis

import (
	"os"
	"sync"
	"time"
)

type LogEntry struct {
	Operation string        `json:"operation"`
	Key       string        `json:"key"`
	Value     string        `json:"value"`
	Ttl       time.Duration `json:"ttl"`
	CreatedAt time.Time     `json:"created_at"`
}

type RedisValue struct {
	value string
	ttl   time.Time
}

type Redis struct {
	data    map[string]RedisValue
	lock    sync.RWMutex
	logfile *os.File
}
