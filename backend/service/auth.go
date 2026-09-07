package service

import (
	"backend/model"
	"backend/repository"
	"context"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *repository.Repository
	jwtKey []byte
}

func NewAuthService(repo *repository.Repository, jwtKey string) *AuthService {
	return &AuthService{
		repo:   repo,
		jwtKey: []byte(jwtKey),
	}
}

func (s *AuthService) Registration(ctx context.Context, name string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	newUser, err := s.repo.InsertUser(ctx, user)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	err = s.repo.InsertCart(ctx, newUser.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("wrong password or email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("wrong password or email")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"cart_id": user.CartId,
	})

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) DataFromJWT(tokenString string) (int, int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return s.jwtKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return 0, 0, err
	}

	if !token.Valid {
		return 0, 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, fmt.Errorf("invalid claims")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("invalid user_id type")
	}

	cartId, ok := claims["cart_id"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("invalid cart_id type")
	}

	return int(userId), int(cartId), nil
}
