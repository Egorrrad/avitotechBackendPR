package service

type Repository interface {
	PullRequestRepo
	TeamRepo
	UserRepo
}

type PullRequestRepo interface {
}

type TeamRepo interface {
}

type UserRepo interface {
}
