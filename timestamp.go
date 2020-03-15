package main

import "time"

type Timestamps struct {
	CreatedAt time.Time  `pg:",default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `pg:",default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `pg:",soft_delete" json:"deleted_at,omitempty"`
}
