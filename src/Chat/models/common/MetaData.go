package common

import "time"

type MetaData struct {
	// Fields for GORM.
	CreatedAt time.Time `json:"toc,omitzero"`
	UpdatedAt time.Time `json:"tou,omitzero"`
}
