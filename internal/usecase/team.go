package usecase

import (
	"context"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

func (s *Service) CreateTeam(ctx context.Context, teamName string, members []domain.TeamMember) (*domain.Team, error) {
	exists, err := s.teams.Exists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTeamAlreadyExists
	}

	domainUsers := make([]domain.User, 0, len(members))
	for _, m := range members {
		domainUsers = append(domainUsers, domain.User{
			UserId:   m.UserId,
			Username: m.Username,
			IsActive: m.IsActive,
			TeamName: teamName,
		})
	}

	team := &domain.Team{
		TeamName: teamName,
		Members:  members,
	}

	if err := s.teams.Create(ctx, team); err != nil {
		return nil, err
	}

	if err := s.users.UpsertBatch(ctx, domainUsers); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.teams.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrTeamNotFound
	}
	return team, nil
}
