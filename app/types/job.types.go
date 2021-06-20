package types

import "time"

type NewJobDTO struct {
	ID        uint
	Backup    uint
	Status    string
	File      string
	CreatedAt time.Time
}

type SmallJob struct {
	Backup      int
	Status      string
	CreatedAt time.Time
}