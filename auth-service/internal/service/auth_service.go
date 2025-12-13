package service

import (
	"errors"
	"time"

	"github.com/dotenv213/aim/auth-service/internal/domain"
	"github.com/dotenv213/aim/auth-service/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo      domain.UserRepository
	jwtSecret     []byte
	tokenDuration time.Duration
}

func NewAuthService(repo domain.UserRepository, secret string) domain.AuthService {
	return &authService{
		userRepo:      repo,
		jwtSecret:     []byte(secret),
		tokenDuration: time.Hour * 72, 
	}
}

func (s *authService) Register(email string, password string) (*domain.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Email:    email,
		Password: hashedPwd,
	}

	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(email string, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials") 
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials") 
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) generateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.tokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}