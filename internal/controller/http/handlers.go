package http

import (
	"github.com/Egorrrad/avitotechBackendPR/internal/usecase"
	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	s usecase.Service
	l logger.Interface
	v *validator.Validate
}
