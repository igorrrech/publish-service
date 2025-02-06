package repo

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	models "github.com/igorrrech/publish-service/authorization/service/models"
)

type TokenRepository struct {
	mx        sync.RWMutex
	cache     map[uint]*models.TokenPair
	secret    string
	accessTTL time.Duration
}

func NewTokenRepository(
	secret string,
	accessTTL time.Duration,
) *TokenRepository {
	return &TokenRepository{
		cache:     make(map[uint]*models.TokenPair),
		secret:    secret,
		accessTTL: 0,
	}
}
func (r TokenRepository) GetAccessByRefresh(refresh models.RefreshToken) (*models.TokenPair, error) {
	r_token, err := models.RefreshToken(refresh).VerifyToken(models.RefreshClaims{}, r.secret)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("refresh verify error:"), err)
	}
	r_claims, ok := r_token.Claims.(*models.RefreshClaims)
	if !ok {
		return nil, fmt.Errorf("refresh claims error")
	}
	r.mx.RLock()
	pair := r.cache[r_claims.UUID]
	r.mx.RUnlock()
	if pair == nil {
		return nil, fmt.Errorf("no such user")
	}
	a_token, err := models.AccesToken(pair.Access).VerifyToken(models.AccessClaims{}, r.secret)
	if (err != nil) && (err != jwt.ErrTokenExpired) {
		return nil, errors.Join(fmt.Errorf("access verify error:"), err)
	}
	a_claims, ok := a_token.Claims.(*models.AccessClaims)
	if !ok {
		return nil, fmt.Errorf("access claims error")
	}
	newPair, err := models.NewTokenPair(
		a_claims.UUID,
		r.accessTTL,
		a_claims.Role,
		r.secret,
	)
	if err != nil {
		return nil, err
	}
	r.mx.Lock()
	defer r.mx.Unlock()
	r.cache[a_claims.UUID] = &newPair
	return &newPair, nil
}
func (r TokenRepository) MakeTokenPair(u models.User) (models.TokenPair, error) {
	pair, err := models.NewTokenPair(
		u.ID,
		r.accessTTL,
		u.Role,
		r.secret,
	)
	r.mx.Lock()
	defer r.mx.Unlock()
	r.cache[u.ID] = &pair
	return pair, err
}
