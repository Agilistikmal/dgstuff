package service

import (
	"context"
	"strconv"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/agilistikmal/dgstuff/internal/pkg"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService struct {
	db           *gorm.DB
	validator    *app.Validator
	tokenService *TokenService
}

func NewAuthService(db *gorm.DB, validator *app.Validator, tokenService *TokenService) *AuthService {
	return &AuthService{db: db, validator: validator, tokenService: NewTokenService()}
}

func (s *AuthService) Me(ctx context.Context, token string) (*model.User, error) {
	claims, err := s.tokenService.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	dataMap, ok := claims["data"].(map[string]any)
	if !ok {
		logrus.Errorf("data map not found in claims: %v", claims)
		return nil, app.NewUnauthorizedError()
	}
	userBytes, err := json.Marshal(dataMap["user"])
	if err != nil {
		logrus.Errorf("failed to marshal user: %v", err)
		return nil, app.NewInternalServerError()
	}
	var user model.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		logrus.Errorf("failed to unmarshal user: %v", err)
		return nil, app.NewInternalServerError()
	}

	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, dto model.UserLoginDTO) (*model.AuthResponse, error) {
	err := s.validator.Validate(dto)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = s.db.Where("email = ?", dto.Email).First(&user).Error
	if err != nil {
		return nil, app.NewUnauthorizedError()
	}

	if !pkg.VerifyPassword(dto.Password, user.Password) {
		return nil, app.NewUnauthorizedError()
	}

	token := s.tokenService.GenerateToken(map[string]any{
		"sub":  strconv.Itoa(user.ID),
		"user": user,
	}, 0)

	return &model.AuthResponse{
		Token: token,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, dto model.UserRegisterDTO) (*model.AuthResponse, error) {
	err := s.validator.Validate(dto)
	if err != nil {
		return nil, err
	}

	var findUser model.User
	s.db.Where("email = ?", dto.Email).First(&findUser)

	if findUser.ID != 0 {
		logrus.Errorf("email already registered: %v", dto.Email)
		return nil, app.NewBadRequestError("email already registered")
	}

	var firstAdmin model.User
	s.db.Where("role = ?", model.UserRoleAdmin).First(&firstAdmin)

	hashedPassword, err := pkg.HashPassword(dto.Password)
	if err != nil {
		logrus.Errorf("failed to hash password: %v", err)
		return nil, err
	}

	user := model.User{
		Email:    dto.Email,
		Password: hashedPassword,
		Role:     model.UserRoleUser,
	}
	if firstAdmin.ID == 0 {
		user.Role = model.UserRoleAdmin
	}

	err = s.db.Save(&user).Error
	if err != nil {
		logrus.Errorf("failed to save user: %v", err)
		return nil, err
	}

	token := s.tokenService.GenerateToken(map[string]any{
		"sub":  strconv.Itoa(user.ID),
		"user": user,
	}, 0)

	return &model.AuthResponse{
		Token: token,
	}, nil
}
