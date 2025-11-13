package model

import "time"

type Good struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"userId"`
	Type      string    `json:"type"` // 'seed', 'bed'
	IDGood    int       `json:"idGood"`
	Quantity  int       `json:"quantity"`
	Cost      int       `json:"cost"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
