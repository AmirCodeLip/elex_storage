package dtos

import (
	"time"

	"github.com/google/uuid"
)

type StorageItemType byte

const (
	Directory StorageItemType = iota
	File
)

type StorageItemDto struct {
	Id        uuid.UUID
	Name      string
	Size      int
	CreatedAt time.Time
	Type      StorageItemType
}
