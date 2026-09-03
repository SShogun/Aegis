package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByProviderSubject(
	ctx context.Context,
	provider string,
	subject string,
) (*User, error) {
	const query = `
		SELECT
			id,
			email,
			name,
			identity_provider,
			provider_subject,
			created_at,
			updated_at
		FROM users
		WHERE identity_provider = $1
		  AND provider_subject = $2
	`

	var user User

	err := r.db.QueryRow(ctx, query, provider, subject).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.IdentityProvider,
		&user.ProviderSubject,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id string,
) (*User, error) {
	const query = `
		SELECT
			id,
			email,
			name,
			identity_provider,
			provider_subject,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.IdentityProvider,
		&user.ProviderSubject,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) Create(
	ctx context.Context,
	user User,
) (*User, error) {
	const query = `
		INSERT INTO users (
			email,
			name,
			identity_provider,
			provider_subject
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			email,
			name,
			identity_provider,
			provider_subject,
			created_at,
			updated_at
	`

	var created User

	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.Name,
		user.IdentityProvider,
		user.ProviderSubject,
	).Scan(
		&created.ID,
		&created.Email,
		&created.Name,
		&created.IdentityProvider,
		&created.ProviderSubject,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
