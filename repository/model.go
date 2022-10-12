package repository

/*
Copyright © 2022 Brad Starkenberg brad@starkenberg.net

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
