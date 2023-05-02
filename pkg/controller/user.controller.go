package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/praadit/dikurium-test/pkg/config"
	"github.com/praadit/dikurium-test/pkg/models"
	"github.com/praadit/dikurium-test/pkg/utils"
	"go.uber.org/zap"
)

func (c *Controller) Signup(ctx context.Context, input *models.SignupInput) (*models.SignupResult, error) {
	genPass, err := utils.GeneratePassword(input.Password)
	if err != nil {
		return nil, err
	}

	err = c.userRepo.Create(c.db, &models.User{
		Email:    input.Email,
		Password: genPass,
	})
	if err != nil {
		c.logger.Info("Failed to create user with details",
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return nil, err
	}

	return &models.SignupResult{
		Email: input.Email,
	}, nil
}
func (c *Controller) Signin(ctx context.Context, input *models.SigninInput) (*models.SigninResult, error) {
	user, err := c.userRepo.GetByEmail(input.Email)
	if err != nil {
		c.logger.Info("Failed to get user with details",
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return nil, err
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return nil, errors.New("Email or Password not match")
	}
	token, err := generateClaimsUser(user)
	if err != nil {
		c.logger.Info("Failed to generate token with details",
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return nil, errors.New("Internal error")
	}

	return &models.SigninResult{
		Email:       user.Email,
		AccessToken: token,
	}, nil
}
func (c *Controller) GetAuthenticatedUser(ctx context.Context) (*models.User, error) {
	claims := ctx.Value("claims").(*jwt.MapClaims)
	if claims == nil {
		return nil, errors.New("Unauthorized")
	}

	claim := *claims

	uuidFin, err := uuid.Parse(claim["user_id"].(string))
	if err != nil {
		return nil, errors.New("Unauthorized")
	}

	user, err := c.userRepo.GetById(uuidFin)
	if err != nil {
		return nil, errors.New("Failed to get authenticated user")
	}

	return user, nil
}

func generateClaimsUser(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"issued_at":  time.Now().Unix(),
		"expires_at": time.Now().Add(12 * time.Hour).Unix(),
		"user_id":    user.ID,
		"email":      user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.Config.ApiSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
