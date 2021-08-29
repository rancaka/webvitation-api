package shared

import (
	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	_uuid := uuid.New()

	out, err := _uuid.MarshalText()
	if err != nil {
		return "", err
	}

	return string(out), err
}
