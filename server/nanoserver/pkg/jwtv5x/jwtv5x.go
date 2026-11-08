package jwtv5x

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Instance *TokenManager
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token expired")
	ErrTokenRevoked  = errors.New("token revoked")
	ErrTokenNotFound = errors.New("token not found")
	// ErrAccountBanned = errors.New("account banned")
	// ErrPlayerBanned  = errors.New("player banned")
)

// GameClaims 自定义声明结构
type GameClaims struct {
	UserId      string `json:"user_id"`
	Name        string `json:"user_name"`
	AccountName string `json:"account_name"`
	Platform    string `json:"platform"`
	ServerId    int32  `json:"server_id"`
	Ad          string `json:"ad,omitempty"`
	jwt.RegisteredClaims
}

// TokenRepo 令牌存储接口
type TokenRepo interface {
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Exist(ctx context.Context, key string) (bool, error)
}

// TokenManager JWT管理器
type TokenManager struct {
	secret        []byte
	tokenDuration time.Duration
	repo          TokenRepo
}

// Config JWT配置
type Config struct {
	Secret        string
	TokenDuration time.Duration
	Repo          TokenRepo
}

// NewTokenManager 创建JWT管理器
func NewTokenManager(cfg Config) {
	tokenMgr := &TokenManager{
		secret:        []byte(cfg.Secret),
		tokenDuration: cfg.TokenDuration,
		repo:          cfg.Repo,
	}
	Instance = tokenMgr
}

// Token 令牌
type Token struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// GenerateToken 生成令牌对
func (tm *TokenManager) GenerateToken(ctx context.Context, userid string, platform, accountName string, serverID int32, adSource string) (string, error) {

	now := time.Now()
	claims := GameClaims{
		UserId: userid,
		// Name:        username,
		AccountName: accountName,
		Platform:    platform,
		ServerId:    serverID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Ad: adSource,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(tm.secret)
	if err != nil {
		return "", err
	}

	if tm.repo != nil {
		key := fmt.Sprintf("token:%s", userid)
		if err = tm.repo.Set(ctx, key, tokenString, tm.tokenDuration); err != nil {
			return "", err
		}
	}

	return tokenString, nil
}

// VerifyToken 验证令牌
func (tm *TokenManager) VerifyToken(ctx context.Context, tokenString string) (*GameClaims, error) {

	gameClaims := &GameClaims{}

	token, err := jwt.ParseWithClaims(tokenString, gameClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return tm.secret, nil
	})

	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ErrExpiredToken
	}

	if err != nil {
		return nil, ErrInvalidToken
	}

	_, ok := token.Claims.(*GameClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if tm.repo != nil {
		key := fmt.Sprintf("token:%s", gameClaims.UserId)
		storedToken, err := tm.repo.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		if storedToken != tokenString {
			return nil, ErrTokenRevoked
		}
	}

	return gameClaims, nil
}
