package common

import "os"

func OpenFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}
	return file, nil
}
