package user

import (
	"fmt"
	"haddibanga/internal/auth"
	"haddibanga/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {
	user := User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.checkPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.Response{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    user.CreatedAt.String(),
	}, nil
}

func (s *service) RefreshToken(refreshToken string) (*dto.Response, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if claims.TokenType != auth.TokenTypeRefresh {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwtService.GenerateAccessToken(claims.UserID, claims.Email, claims.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &dto.Response{
		AccessToken: accessToken,
	}, nil
}
