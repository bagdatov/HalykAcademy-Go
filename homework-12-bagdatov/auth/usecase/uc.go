package uc

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"web/domain"

	"github.com/dgrijalva/jwt-go"
)

type authUseCase struct {
	accessSecret  string
	refreshSecret string
	cache         domain.CaсheStore
	db            domain.Repository
}

func NewAuthUseCase(db domain.Repository, casheStore domain.CaсheStore) *authUseCase {
	return &authUseCase{
		accessSecret:  getRandomSecret(),
		refreshSecret: getRandomSecret(),
		cache:         casheStore,
		db:            db,
	}
}

func getRandomSecret() string {
	var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	rand.Seed(time.Now().UnixNano())

	str := make([]rune, 32)

	for i := range str {
		str[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(str)
}

func (a *authUseCase) Login(username, password string) (*domain.User, error) {
	return a.db.FindUser(context.Background(), username, password)
}

func (a *authUseCase) ExtractToken(AuthorizationHeader string) (string, error) {
	// header := r.Header.Get("Authorization")
	if AuthorizationHeader == "" {
		return "", domain.ErrNoHeader
	}

	parsedHeader := strings.Split(AuthorizationHeader, " ")
	if len(parsedHeader) != 2 || parsedHeader[0] != "Bearer" {
		return "", domain.ErrInvalidHeader
	}

	return parsedHeader[1], nil
}

func (a *authUseCase) ParseToken(token string, isAccess bool) (*domain.User, error) {

	JWTToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to extract token metadata, unexpected signing method: %v", token.Header["alg"])
		}
		if isAccess {
			return []byte(a.accessSecret), nil
		}
		return []byte(a.refreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := JWTToken.Claims.(jwt.MapClaims)

	if ok && JWTToken.Valid {

		var userID float64

		userID, ok = claims["id"].(float64)
		if !ok {
			return nil, domain.ErrInvalidToken
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, domain.ErrInvalidToken
		}

		userName, ok := claims["name"].(string)
		if !ok {
			return nil, domain.ErrInvalidToken
		}

		expiredTime := time.Unix(int64(exp), 0)

		if time.Now().After(expiredTime) {
			return nil, domain.ErrExpiredToken
		}
		return &domain.User{
			ID:       int(userID),
			Username: userName,
		}, nil
	}

	return nil, domain.ErrInvalidToken
}

// GenerateAndSendTokens return access token and refresh token in that order.
func (a *authUseCase) GenerateAndSendTokens(u *domain.User) (string, string, error) {

	accessTokenExp := time.Now().Add(a.cache.GetAccessTokenTTL()).Unix()

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["id"] = u.ID
	accessTokenClaims["name"] = u.Username
	accessTokenClaims["iat"] = time.Now().Unix()
	accessTokenClaims["exp"] = accessTokenExp
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessSignedToken, err := accessToken.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", "", domain.ErrTokenNotCreated
	}

	refreshTokenExp := time.Now().Add(a.cache.GetRefreshTokenTTL()).Unix()
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["id"] = u.ID
	refreshTokenClaims["name"] = u.Username
	refreshTokenClaims["iat"] = time.Now().Unix()
	refreshTokenClaims["exp"] = refreshTokenExp
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshSignedToken, err := refreshToken.SignedString([]byte(a.refreshSecret))
	if err != nil {
		return "", "", domain.ErrTokenNotCreated
	}

	if err := a.cache.InsertToken(int64(u.ID), refreshSignedToken); err != nil {
		fmt.Println(err)
		return "", "", domain.ErrTokenNotCreated
	}
	return accessSignedToken, refreshSignedToken, nil
}

// UpdateToken returns access token and refresh token (in that order).
func (a *authUseCase) UpdateToken(refreshToken string) (string, string, error) {

	//checking update token
	user, err := a.ParseToken(refreshToken, false)
	if err != nil {
		return "", "", err
	}

	// find token in Redis
	ok := a.cache.FindToken(int64(user.ID), refreshToken)
	if !ok {
		return "", "", domain.ErrInvalidToken
	}

	return a.GenerateAndSendTokens(user)
}
