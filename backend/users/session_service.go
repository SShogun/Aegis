package users

import (
	"context"
	"time"
)

const SessionDuration = 24 * time.Hour

type SessionService struct {
	repository *SessionRepository
}

func NewSessionService(repository *SessionRepository) *SessionService {
	return &SessionService{
		repository: repository,
	}
}

func (s *SessionService) Create(
	ctx context.Context,
	userID string,
) (*Session, error) {
	return s.repository.Create(
		ctx,
		userID,
		SessionDuration,
	)
}

func (s *SessionService) FindValid(
	ctx context.Context,
	sessionID string,
) (*Session, error) {
	return s.repository.FindValid(ctx, sessionID)
}

func (s *SessionService) Delete(
	ctx context.Context,
	sessionID string,
) error {
	return s.repository.Delete(ctx, sessionID)
}
