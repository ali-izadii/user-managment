package handler

import (
	"github.com/docker/docker/daemon/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	validator *validator.Validate
	logger    logger.Logger
}

func NewUserHandler(validator *validator.Validate, logger logger.Logger) *UserHandler {
	return &UserHandler{
		validator: validator,
		logger:    logger,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
}

func (h *UserHandler) GetProfiles(c *gin.Context) {
}
