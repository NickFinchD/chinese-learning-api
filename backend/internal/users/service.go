package users

import (
	"context"
	"errors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AdminList(ctx context.Context, request AdminListRequest) ([]User, int, error) {
	return s.repository.AdminList(ctx, request)
}

// SetAdmin refuses to revoke the last remaining admin — otherwise a
// mis-click here locks everyone out of the admin panel, with no in-app way
// back in (only the promote-admin CLI could recover it).
func (s *Service) SetAdmin(ctx context.Context, id int64, isAdmin bool) (*User, error) {

	if !isAdmin {

		target, err := s.repository.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if target.IsAdmin {

			count, err := s.repository.CountAdmins(ctx)
			if err != nil {
				return nil, err
			}

			if count <= 1 {
				return nil, errors.New("cannot revoke the last remaining admin")
			}
		}
	}

	return s.repository.SetAdmin(ctx, id, isAdmin)
}
