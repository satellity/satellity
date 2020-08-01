package models

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Session contains user's current login infomation
type Session struct {
	SessionID    string
	UserID       string
	PrivateKey   string
	ClientPublic string
	CreatedAt    time.Time

	PublicKey string
}

var sessionColumns = []string{"session_id", "user_id", "private_key", "client_public", "created_at"}

func (s *Session) values() []interface{} {
	return []interface{}{s.SessionID, s.UserID, s.PrivateKey, s.ClientPublic, s.CreatedAt}
}

func sessionFromRows(row durable.Row) (*Session, error) {
	var s Session
	err := row.Scan(&s.SessionID, &s.UserID, &s.PrivateKey, &s.ClientPublic, &s.CreatedAt)
	return &s, err
}

// CreateSession create a new user session
func CreateSession(ctx context.Context, identity, password, public string) (*User, error) {
	data, err := base64.RawURLEncoding.DecodeString(public)
	if err != nil {
		return nil, session.BadDataError(ctx)
	}
	pub, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, session.BadDataError(ctx)
	}
	switch pub.(type) {
	case ed25519.PublicKey:
	default:
		return nil, session.BadDataError(ctx)
	}

	user, err := ReadUserByUsernameOrEmail(ctx, identity)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, session.IdentityNonExistError(ctx)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword.String), []byte(password)); err != nil {
		return nil, session.InvalidPasswordError(ctx)
	}

	err = session.Database(ctx).RunInTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sessions WHERE session_id IN (SELECT session_id FROM sessions WHERE user_id=$1 ORDER BY created_at DESC OFFSET 5)", user.UserID)
		if err != nil {
			return err
		}
		s, err := user.addSession(ctx, tx, public)
		if err != nil {
			return err
		}
		user.Session = s
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func (user *User) addSession(ctx context.Context, tx *sql.Tx, public string) (*Session, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	s := &Session{
		SessionID:    uuid.Must(uuid.NewV4()).String(),
		UserID:       user.UserID,
		PrivateKey:   base64.RawURLEncoding.EncodeToString(priv),
		ClientPublic: public,
		CreatedAt:    time.Now(),
		PublicKey:    base64.RawURLEncoding.EncodeToString(pub),
	}

	columns, positions := durable.PrepareColumnsWithParams(sessionColumns)
	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf("INSERT INTO sessions(%s) VALUES(%s)", columns, positions))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, s.values()...)
	return s, nil
}

func readSession(ctx context.Context, tx *sql.Tx, uid, sid string) (*Session, error) {
	if id, _ := uuid.FromString(uid); id.String() == uuid.Nil.String() {
		return nil, nil
	}
	if id, _ := uuid.FromString(sid); id.String() == uuid.Nil.String() {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM sessions WHERE user_id=$1 AND session_id=$2", strings.Join(sessionColumns, ",")), uid, sid)
	s, err := sessionFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func PrivateKeyToCurve25519(curve25519Private *[32]byte, privateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(privateKey[:32])
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(curve25519Private[:], digest)
}

func PublicKeyToCurve25519(curve25519Public *[32]byte, publicKey ed25519.PublicKey) error {
	var k [32]byte
	copy(k[:], publicKey[:])
	var A ExtendedGroupElement
	if !A.FromBytes(&k) {
		return fmt.Errorf("Invalid public key %x", publicKey)
	}

	// A.Z = 1 as a postcondition of FromBytes.

	var x FieldElement
	edwardsToMontgomeryX(&x, &A.Y)
	FeToBytes(curve25519Public, &x)
	return nil
}

func edwardsToMontgomeryX(outX, y *FieldElement) {
	// We only need the x-coordinate of the curve25519 point, which I'll
	// call u. The isomorphism is u=(y+1)/(1-y), since y=Y/Z, this gives
	// u=(Y+Z)/(Z-Y). We know that Z=1, thus u=(Y+1)/(1-Y).
	var oneMinusY FieldElement
	FeOne(&oneMinusY)
	FeSub(&oneMinusY, &oneMinusY, y)
	FeInvert(&oneMinusY, &oneMinusY)

	FeOne(outX)
	FeAdd(outX, outX, y)

	FeMul(outX, outX, &oneMinusY)
}
