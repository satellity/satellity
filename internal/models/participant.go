package models

import "time"

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
const (
	ParticipantRoleOwner  = "OWNER"
	ParticipantRoleAdmin  = "ADMIN"
	ParticipantRoleVIP    = "VIP"
	ParticipantRoleMember = "MEMBER"
)

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
