package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"api/src/config"

	"api/src/common/security"
	"api/src/modules/auth/dto"
	"api/src/modules/users"
	usersDto "api/src/modules/users/dto"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(d dto.LoginDto) (*dto.TokenResponseDto, error)
	Register(d usersDto.CreateUserDto) (*dto.TokenResponseDto, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	Logout(tokenString string) error
	IsBlacklisted(tokenString string) bool
}

type authService struct {
	userService users.UserService
}

func NewAuthService(userService users.UserService) AuthService {
	return &authService{userService}
}

func (s *authService) Login(d dto.LoginDto) (*dto.TokenResponseDto, error) {
	user, err := s.userService.FindByEmail(d.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !security.ComparePasswords(d.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateToken(user.ID, user.Role, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, user.Role, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Register(d usersDto.CreateUserDto) (*dto.TokenResponseDto, error) {
	user, err := s.userService.Create(d)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateToken(user.ID, user.Role, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, user.Role, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) generateToken(userID uint, role string, duration time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_change_me"
	}

	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(duration).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	if s.IsBlacklisted(tokenString) {
		return nil, errors.New("token is blacklisted")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_change_me"
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
}

func (s *authService) Logout(tokenString string) error {
	if config.RedisClient == nil {
		return nil
	}

	return config.RedisClient.Set(context.Background(), "blacklist:"+tokenString, "true", 24*time.Hour).Err()
}

func (s *authService) IsBlacklisted(tokenString string) bool {
	if config.RedisClient == nil {
		return false
	}

	val, err := config.RedisClient.Get(context.Background(), "blacklist:"+tokenString).Result()
	return err == nil && val == "true"
}
