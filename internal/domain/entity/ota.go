package entity

import "time"

type OTA struct {
	ID           string    `json:"id" db:"id"`
	AppID        string    `json:"app_id" db:"app_id"`
	VersionName  string    `json:"version_name" db:"version_name"`
	VersionCode  int       `json:"version_code" db:"version_code"`
	ReleaseNotes string    `json:"release_notes" db:"release_notes"`
	URL          string    `json:"url" db:"url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
