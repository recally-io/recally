package cache

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"recally/internal/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CacheTestSuite struct {
	suite.Suite
	*DbCache
	ctx       context.Context
	container *postgres.PostgresContainer
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *CacheTestSuite) SetupSuite() {
	// Create a new context
	ctx := context.Background()
	s.ctx = ctx
	// set up pg container
	pool, err := s.initPostgres(ctx)
	if err != nil {
		log.Fatalf("could not setup postgres container: %v", err)
	}
	// Create a new cache instance
	s.DbCache = NewDBCache(pool)
}

func (s *CacheTestSuite) TearDownSuite() {
	// Clean up the container
	if err := s.container.Terminate(s.ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func (s *CacheTestSuite) initPostgres(ctx context.Context) (*db.Pool, error) {
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithInitScripts(filepath.Join("testdata", "000001_new_cache_table.up.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	s.container = postgresContainer

	// Create a new pool
	ep, err := postgresContainer.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	eps := strings.Split(ep, ":")

	databaseUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", eps[0], eps[1], dbUser, dbPassword, dbName)

	pool, err := db.NewPool(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (s *CacheTestSuite) TestCacheString() {
	// Test cache
	ctx := context.Background()
	domain := "string"
	key := "key"
	value := "string data"
	cacheKey := NewCacheKey(domain, key)
	s.SetWithContext(ctx, cacheKey, value, 1*time.Hour)

	val, ok := Get[string](ctx, s.DbCache, cacheKey)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), value, *val)
}

func (s *CacheTestSuite) TestCacheStructData() {
	// Test cache
	ctx := context.Background()
	domain := "users"
	key := "1"
	value := User{ID: 1, Name: "John"}
	cacheKey := NewCacheKey(domain, key)
	s.SetWithContext(ctx, cacheKey, value, 1*time.Hour)

	user, ok := Get[User](ctx, s.DbCache, cacheKey)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), value, *user)
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
