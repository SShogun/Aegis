package users

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	userID string,
	duration time.Duration,
) (*Session, error) {
	expiresAt := time.Now().Add(duration)

	const query = `
		INSERT INTO sessions (
			id,
			user_id,
			expires_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2
		)
		RETURNING
			id,
			user_id,
			created_at,
			expires_at
	`

	var session Session

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		expiresAt,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) FindValid(
	ctx context.Context,
	sessionID string,
) (*Session, error) {
	const query = `
		SELECT
			id,
			user_id,
			created_at,
			expires_at
		FROM sessions
		WHERE id = $1
		  AND expires_at > NOW()
	`

	var session Session

	err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) Delete(
	ctx context.Context,
	sessionID string,
) error {
	const query = `
		DELETE FROM sessions
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, sessionID)

	return err
}
