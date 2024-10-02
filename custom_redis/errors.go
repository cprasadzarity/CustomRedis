package custom_redis

import (
	"CustomRedis/common"
)

const (
	KEY_NOT_FOUND_ERROR = 1001
)

func KeyNotFoundError() error {
	return &common.CustomError{
		Message: "Key Not Found Error",
		Code:    KEY_NOT_FOUND_ERROR,
		Details: "Could not find the key in redis",
	}
}
