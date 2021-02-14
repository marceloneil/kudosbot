package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Kudos struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Timestamp  time.Time `db:"timestamp" json:"timestamp"`
	SenderID   string    `db:"sender_id" json:"sender_id"`
	ReceiverID string    `db:"receiver_id" json:"receiver_id"`
	EventID    string    `db:"event_id" json:"event_id"`
	Reaction   string    `db:"reaction" json:"reaction"`
}
