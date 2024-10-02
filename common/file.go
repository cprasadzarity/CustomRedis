package common

import "os"

func OpenFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	return file, nil
}
