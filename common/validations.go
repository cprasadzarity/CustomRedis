package common

import "errors"

func ValidateRequiredKeys(data map[string]string, keys ...string) error {
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			return errors.New("Required key " + key + " not found")
		}
	}
	return nil
}
