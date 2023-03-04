package utils

import (
	"github.com/obgnail/plugin-platform/platform/service/utils/uuid"
)

func NewInstanceUUID() string {
	return uuid.UUID()
}
