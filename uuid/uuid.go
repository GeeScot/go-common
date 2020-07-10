package uuid

import guuid "github.com/google/uuid"

func Token() string {
	return guuid.New().String()
}
