package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type RequestLog struct {
	ID         uuid.UUID `gorm:"type:uuid;"`
	RawRequest string
	CreatedAt  time.Time
}

type FailedRequestLog struct {
	ID         uuid.UUID `gorm:"type:uuid;"`
	RawRequest string
	CreatedAt  time.Time
}
