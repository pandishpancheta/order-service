package model

import "github.com/gofrs/uuid/v5"

type Status string

const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Failed    Status = "failed"
)

type Order struct {
	CreatedAt string    `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ListingID uuid.UUID `json:"listing_id"`
	Status    Status    `json:"status"`
	TokenURI  string    `json:"token_uri"`
}
