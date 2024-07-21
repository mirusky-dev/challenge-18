package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/mirusky-dev/challenge-18/core"
	"github.com/mirusky-dev/challenge-18/core/env"
	"github.com/mirusky-dev/challenge-18/repositories"
)

type ITokenService interface {
	Issue(ctx context.Context, userID string) (string, string, time.Time, time.Time, *core.Exception)
	Refresh(ctx context.Context, refreshToken string) (string, string, time.Time, time.Time, *core.Exception)
	Revoke(ctx context.Context, tokenJTI string, expiresAt time.Time) *core.Exception
	IsRevoked(ctx context.Context, tokenJTI string) *core.Exception
}

type jwtTokenService struct {
	config env.Config

	refreshTokenStorage    fiber.Storage
	tokenRevokationStorage fiber.Storage
	userRepository         repositories.IUserRepository
}

func NewTokenService(
	config env.Config,
	refreshTokenStorage fiber.Storage,
	tokenRevokationStorage fiber.Storage,
	userRepository repositories.IUserRepository,
) ITokenService {
	return &jwtTokenService{
		config: config,

		refreshTokenStorage:    refreshTokenStorage,
		tokenRevokationStorage: tokenRevokationStorage,
		userRepository:         userRepository,
	}
}

// Issue creates a new JWT token to the given userID
//
//	ref is a userID (TokenRef) or a userID;Signature (RefreshTokenRef)
func (svc *jwtTokenService) Issue(ctx context.Context, ref string) (string, string, time.Time, time.Time, *core.Exception) {

	parts := strings.Split(ref, ";")

	userID := parts[0]

	user, errUser := svc.userRepository.GetByID(ctx, userID)
	if errUser != nil {
		return "", "", time.Time{}, time.Time{}, errUser
	}

	// TODO?: Should we invalidate all refresh tokens when user changes password?
	// If parts == 2, userIDRef is a refreshRef, so we need to also check the signature
	// Signature changes when user changes password, so all refresh token are revoked
	if len(parts) == 2 {
		if !strings.EqualFold(parts[1], user.Signature) {
			return "", "", time.Time{}, time.Time{}, core.BadRequest(core.WithMessage("invalid or expired refresh token"))
		}
	}

	now := time.Now()

	// TODO?: Get expiration from environment variables or use default?
	expiresAt := now.Add(time.Minute * 5)

	claims := core.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Role: user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(svc.config.JWTSecret))
	if err != nil {
		return "", "", time.Time{}, time.Time{}, core.Unexpected(core.WithError(err))
	}

	refreshToken := uuid.New().String()
	expiresIn := time.Hour
	refreshExpiresAt := now.Add(expiresIn)

	refreshTokenRef := fmt.Sprintf("%s;%s", user.ID, user.Signature)

	if err := svc.refreshTokenStorage.Set(refreshToken, []byte(refreshTokenRef), expiresIn); err != nil {
		return "", "", time.Time{}, time.Time{}, core.Unexpected(core.WithError(err))
	}

	return tokenStr, refreshToken, expiresAt, refreshExpiresAt, nil
}

// Refresh check if refresh token still valid and then generate a new one
func (svc *jwtTokenService) Refresh(ctx context.Context, refreshToken string) (string, string, time.Time, time.Time, *core.Exception) {

	refreshRef, err := svc.refreshTokenStorage.Get(refreshToken)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, core.Unexpected(core.WithError(err))
	}

	if refreshRef == nil {
		return "", "", time.Time{}, time.Time{}, core.BadRequest(core.WithMessage("invalid or expired refresh token"))
	}

	if err := svc.refreshTokenStorage.Delete(refreshToken); err != nil {
		return "", "", time.Time{}, time.Time{}, core.Unexpected(core.WithError(err))
	}

	return svc.Issue(ctx, string(refreshRef))
}

// Revoke is responsable to mark JWT token as revoked
func (svc *jwtTokenService) Revoke(ctx context.Context, tokenJTI string, tokenExpiresAt time.Time) *core.Exception {

	err := svc.tokenRevokationStorage.Set(tokenJTI, []byte(tokenJTI), time.Until(tokenExpiresAt))
	if err != nil {
		return core.Unexpected(core.WithError(err))
	}

	return nil
}

// IsRevoked ...
func (svc *jwtTokenService) IsRevoked(ctx context.Context, tokenJTI string) *core.Exception {

	value, err := svc.tokenRevokationStorage.Get(tokenJTI)
	if err != nil {
		return core.Unexpected(core.WithError(err))
	} else if value != nil {
		return core.Forbidden(core.WithMessage("Token has been revoked"))
	}

	return nil
}
