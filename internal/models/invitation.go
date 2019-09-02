package models

import (
	"satellity/internal/durable"
	"time"
)

const invitationsDDL = ``

// Invitation is a way to invate user to group for free
type Invitation struct {
	InvitationID string
	GroupID      string
	Email        string
	Code         string
	SentAt       time.Time
	CreatedAt    time.Time
}

var invitationColumns = []string{"invitation_id", "group_id", "email", "code", "sent_at", "created_at"}

func (i *Invitation) values() []interface{} {
	return []interface{}{i.InvitationID, i.GroupID, i.Email, i.Code, i.SentAt, i.CreatedAt}
}

func invitationFromRows(row durable.Row) (*Invitation, error) {
	var i Invitation
	err := row.Scan(&i.InvitationID, &i.GroupID, &i.Email, &i.Code, &i.SentAt, &i.CreatedAt)
	return &i, err
}
