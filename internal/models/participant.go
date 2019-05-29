package models

import (
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"time"
)

const participantsDDL = `
CREATE TABLE IF NOT EXISTS participants (
	group_id               VARCHAR(36) NOT NULL REFERENCES groups ON DELETE CASCADE,
	user_id                VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	role                   VARCHAR(128) NOT NULL,
	created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	PRIMARY KEY (group_id, user_id)
);
`

// Roles of the participant
const (
	ParticipantRoleOwner  = "OWNER"
	ParticipantRoleAdmin  = "ADMIN"
	ParticipantRoleVIP    = "VIP"
	ParticipantRoleMember = "MEMBER"
)

// Participant represents the struct of a group member
type Participant struct {
	GroupID   string
	UserID    string
	Role      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

var participantColumns = []string{"group_id", "user_id", "role", "created_at", "updated_at"}

func (p *Participant) values() []interface{} {
	return []interface{}{p.GroupID, p.UserID, p.Role, p.CreatedAt, p.UpdateAt}
}

func participantFromRow(row durable.Row) (*Participant, error) {
	var p Participant
	err := row.Scan(&p.GroupID, &p.UserID, &p.Role, &p.CreatedAt, &p.UpdateAt)
	return &p, err
}

func createParticipant(mcxt *Context, tx *sql.Tx, groupID, userID, role string) (*Participant, error) {
	ctx := mcxt.context
	t := time.Now()
	p := &Participant{
		GroupID:   groupID,
		UserID:    userID,
		Role:      role,
		CreatedAt: t,
		UpdateAt:  t,
	}

	columns, params := durable.PrepareColumnsWithValues(participantColumns)
	query := fmt.Sprintf("INSERT INTO participants(%s) VALUES (%s)", columns, params)
	_, err := tx.ExecContext(ctx, query, p.values()...)
	return p, err
}
