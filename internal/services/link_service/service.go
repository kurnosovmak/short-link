package link_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/kurnosovmak/short-link/pkg/cachemap"
	"github.com/kurnosovmak/short-link/pkg/logging"
	"github.com/kurnosovmak/short-link/pkg/redis"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/thanhpk/randstr"
	"time"
)

type service struct {
	Kdb    redis.KVContract
	Cache  cachemap.CacheContract
	Logger logging.Logger
}

func NewService(kdb redis.KVContract, cache cachemap.CacheContract, logger logging.Logger) LinkService {
	return &service{
		Kdb:    kdb,
		Cache:  cache,
		Logger: logger,
	}
}

type LinkService interface {
	RedirectLink(ctx context.Context, dto RedirectLinkDTO) (string, error)
	CreateLink(ctx context.Context, dto CreateLinkDTO) (string, error)
}

func (s *service) RedirectLink(ctx context.Context, dto RedirectLinkDTO) (string, error) {
	var link string
	cacheLink := s.Cache.Get(dto.Key)
	if cacheLink != nil {
		return fmt.Sprintf("%v", cacheLink), nil
	}

	link, err := s.Kdb.Get(ctx, dto.Key)
	if err != nil {
		return "", err
	}

	if link == "" {
		return "", errors.New("not fund link")
	}
	err = s.Cache.Set(dto.Key, link, time.Minute*3)
	if err != nil {
		s.Logger.Error(err)
		return "", err
	}

	return link, nil
}
func (s *service) CreateLink(ctx context.Context, dto CreateLinkDTO) (string, error) {
	key := randstr.Hex(16)
	_, err := s.Kdb.Set(ctx, key, dto.Link, redis2.KeepTTL)
	if err != nil {
		return "", err
	}

	return key, nil
}
