package users

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("Invalid email or password")

type Service struct {
	repo      *Repository
	jwtSecret []byte
	jwtExpHrs int
}

func NewService(repo *Repository, jwtSecret string, jwtExpHrs int) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		jwtExpHrs: jwtExpHrs,
	}
}
func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// hashear password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         "customer",
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	token, err := s.generateToken(user)
	if err != nil {

		return nil, err
	}
	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil

}
func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {

		return nil, ErrInvalidCredentials

	}
	// comparar contraseña contra el hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := s.generateToken(user)
	if err != nil {

		return nil, err
	}
	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil

}
func (s *Service) FindByID(id string) (*User, error) {

	return s.repo.FindByID(id)
}
func (s *Service) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Duration(s.jwtExpHrs) * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
func (s *Service) ValidateToken(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token inválido")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims inválidos")
	}
	return &claims, nil
}
