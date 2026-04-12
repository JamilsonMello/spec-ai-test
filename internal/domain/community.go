package domain

import (
	"time"
)

type Community struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (c *Community) IsValidName() bool {
	return len(c.Name) > 0 && len(c.Name) <= 100
}

func (c *Community) IsValidDescription() bool {
	return len(c.Description) > 0 && len(c.Description) <= 500
}
