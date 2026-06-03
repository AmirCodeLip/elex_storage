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

type DirectoryDto struct {
	Id       uuid.UUID
	Name     string
	Size     int
	ParentId uuid.UUID
}

type FileDto struct {
	Id            uuid.UUID
	Name          string
	FileExtension string
	Size          int
	DirectoryId   uuid.UUID
	CreatedAt     time.Time
	UpdateAt      time.Time
}
