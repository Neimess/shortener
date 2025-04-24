package url

import (
	"context"
	"errors"
	"time"

	cachepkg "github.com/Neimess/shortener/internal/infrastructure/cache"
	"github.com/Neimess/shortener/internal/infrastructure/db"
	genpkg "github.com/Neimess/shortener/internal/infrastructure/generator"
	"github.com/Neimess/shortener/internal/model/url"
	"github.com/google/uuid"
)

type serviceImpl struct {
	repos     db.RepositorySet
	cache     cachepkg.ExpirableCache
	generator *genpkg.CharsetGenerator
}

func NewService(
	repos db.RepositorySet,
	cache cachepkg.ExpirableCache,
	gen *genpkg.CharsetGenerator,

) Service {
	return &serviceImpl{repos: repos, cache: cache, generator: gen}
}

func (s *serviceImpl) Shorten(original string) (string, error) {
	if original == "" {
		return "", errors.New("original URL is empty")
	}
	code := s.generator.Generate(6)
	u := &url.URL{
		ID:        uuid.NewString(),
		Original:  original,
		ShortCode: code,
		CreatedAt: time.Now(),
	}
	ctx := context.Background()
	if err := s.repos.URL.Save(ctx, u); err != nil {
		return "", err
	}
	if err := s.cache.SetWithTTL(ctx, code, original, 5*time.Minute); err != nil {
		return "", err
	}
	return code, nil
}

func (s *serviceImpl) Resolve(code string) (string, error) {
	ctx := context.Background()
	if orig, err := s.cache.Get(ctx, code); err == nil && orig != "" {
		return orig, nil
	}
	u, err := s.repos.URL.FindByCode(ctx, code)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("not found")
	}
	_ = s.cache.SetWithTTL(ctx, code, u.Original, 5*time.Minute)
	return u.Original, nil
}
