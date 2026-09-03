package users

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindOrCreateFromClaims(
	ctx context.Context,
	provider string,
	subject string,
	email string,
	name string,
) (*User, error) {
	user, err := s.repository.FindByProviderSubject(
		ctx,
		provider,
		subject,
	)

	if err == nil {
		return user, nil
	}

	if err != ErrNotFound {
		return nil, err
	}

	return s.repository.Create(ctx, User{
		Email:            email,
		Name:             name,
		IdentityProvider: provider,
		ProviderSubject:  subject,
	})
}

func (s *Service) FindByID(ctx context.Context, userID string) (*User, error) {
	return s.repository.FindByID(ctx, userID)
}
