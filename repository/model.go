package repository

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string
	HostName string
	Videos   []Video
}

type Video struct {
	gorm.Model
	UserID      uint   `gorm:"uniqueIndex:idx_videos"`
	PlexGuid    string `gorm:"uniqueIndex:idx_videos"`
	Title       string
	Type        string
	Year        int
	AspectRatio float32 // > 1.5 = widescreen
	Files       []VideoFile
}

type VideoFile struct {
	gorm.Model
	VideoID    uint `gorm:"index:idx_files"`
	FileName   string
	Size       int
	Resolution int
}
