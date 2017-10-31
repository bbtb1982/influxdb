package organizations

import (
	"context"
	"fmt"

	"github.com/influxdata/chronograf"
)

// Ensure UsersStore implements chronograf.UsersStore.
var _ chronograf.UsersStore = &UsersStore{}

// UsersStore uses bolt to store and retrieve users
type UsersStore struct {
	organization string
	store        chronograf.UsersStore
}

func NewUsersStore(s chronograf.UsersStore, org string) *UsersStore {
	return &UsersStore{
		store:        s,
		organization: org,
	}
}

// validOrganizationRoles ensures that each User Role has both an associated Organization and a Name
func validOrganizationRoles(orgID string, u *chronograf.User) error {
	if u == nil || u.Roles == nil {
		return nil
	}
	for _, r := range u.Roles {
		if r.Organization == "" {
			return fmt.Errorf("user role must have an Organization")
		}
		if r.Organization != orgID {
			return fmt.Errorf("organizationID %s does not match %s", r.Organization, orgID)
		}
		if r.Name == "" {
			return fmt.Errorf("user role must have a Name")
		}
	}
	return nil
}

// Get searches the UsersStore for user with name
func (s *UsersStore) Get(ctx context.Context, q chronograf.UserQuery) (*chronograf.User, error) {
	err := validOrganization(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := s.store.Get(ctx, q)
	if err != nil {
		return nil, err
	}

	// filter Roles that are not scoped to this Organization
	roles := usr.Roles[:0]
	for _, r := range usr.Roles {
		if r.Organization == s.organization {
			roles = append(roles, r)
		}
	}
	if len(roles) == 0 {
		// This means that the user has no roles in an organization
		// TODO: should we return the user without any roles or ErrUserNotFound?
		return nil, chronograf.ErrUserNotFound
	}
	usr.Roles = roles
	return usr, nil
}

// Add a new User to the UsersStore.
func (s *UsersStore) Add(ctx context.Context, u *chronograf.User) (*chronograf.User, error) {
	err := validOrganization(ctx)
	if err != nil {
		return nil, err
	}

	if err := validOrganizationRoles(s.organization, u); err != nil {
		return nil, err
	}
	usr, err := s.store.Get(ctx, chronograf.UserQuery{
		Name:     &u.Name,
		Provider: &u.Provider,
		Scheme:   &u.Scheme,
	})
	if err != nil && err != chronograf.ErrUserNotFound {
		return nil, err
	}
	if err == chronograf.ErrUserNotFound {
		return s.store.Add(ctx, u)
	}
	usr.Roles = append(usr.Roles, u.Roles...)
	if err := s.store.Update(ctx, usr); err != nil {
		return nil, err
	}
	u.ID = usr.ID
	return u, nil
}

// Delete a user from the UsersStore
func (s *UsersStore) Delete(ctx context.Context, usr *chronograf.User) error {
	err := validOrganization(ctx)
	if err != nil {
		return err
	}
	u, err := s.store.Get(ctx, chronograf.UserQuery{ID: &usr.ID})
	if err != nil {
		return err
	}
	// delete Roles that are not scoped to this Organization
	roles := u.Roles[:0]
	for _, r := range u.Roles {
		if r.Organization != s.organization {
			roles = append(roles, r)
		}
	}
	u.Roles = roles
	return s.store.Update(ctx, u)
}

// Update a user
func (s *UsersStore) Update(ctx context.Context, usr *chronograf.User) error {
	err := validOrganization(ctx)
	if err != nil {
		return err
	}

	if err := validOrganizationRoles(s.organization, usr); err != nil {
		return err
	}
	u, err := s.store.Get(ctx, chronograf.UserQuery{ID: &usr.ID})
	if err != nil {
		return err
	}
	// filter Roles that are not scoped to this Organization
	roles := u.Roles[:0]
	for _, r := range u.Roles {
		if r.Organization != s.organization {
			roles = append(roles, r)
		}
	}

	// recombine roles from usr, by replacing the roles of the user
	// within the current Organization
	u.Roles = append(roles, usr.Roles...)

	return s.store.Update(ctx, u)
}

// All returns all users
func (s *UsersStore) All(ctx context.Context) ([]chronograf.User, error) {
	err := validOrganization(ctx)
	if err != nil {
		return nil, err
	}

	usrs, err := s.store.All(ctx)
	if err != nil {
		return nil, err
	}

	// Filtering users that have no associtation to an organization
	us := usrs[:0]
	for _, usr := range usrs {
		roles := usr.Roles[:0]
		for _, r := range usr.Roles {
			if r.Organization == s.organization {
				roles = append(roles, r)
			}
		}
		if len(roles) != 0 {
			// Only add users if they have a role in the associated organization
			usr.Roles = roles
			us = append(us, usr)
		}
	}

	return us, nil
}
