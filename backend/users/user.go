package users

import "time"

// User represents an authenticated Aegis user.
//
// The identity provider establishes who the user is.
// Aegis keeps this internal representation so application
// logic does not depend directly on the external provider.
type User struct {
	ID               string
	Email            string
	Name             string
	IdentityProvider string
	ProviderSubject  string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
