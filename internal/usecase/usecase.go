package usecase

type (
	Postgres interface {
	}

	TeamRepo interface {
	}

	UserRepo interface {
	}
	PullRequestRepo interface {
	}
)

type Service struct {
	teams TeamRepo
	users UserRepo
	pr    PullRequestRepo
}

func NewService(team TeamRepo, users UserRepo, pr PullRequestRepo) *Service {
	return &Service{
		teams: team,
		users: users,
		pr:    pr,
	}
}

type Repository interface {
	PullRequestRepo
	TeamRepo
	UserRepo
}

type Service2 struct {
	repo Repository
}
