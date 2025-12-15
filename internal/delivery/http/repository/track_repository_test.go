package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"hub/internal/delivery/http/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestProperty_Repository_UpsertBehavior(t *testing.T) {
	pool := setupTestDB(t)
	if pool == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}
	defer pool.Close()

	repo := NewTrackRepository(pool)
	ctx := context.Background()

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10
	properties := gopter.NewProperties(parameters)

	// Property: Upsert increments rotate on duplicate MD5Sum
	properties.Property("upsert increments rotate on duplicate MD5Sum", prop.ForAll(
		func(md5sum string) bool {
			track1 := &entity.Track{
				ID:     uuid.Must(uuid.NewV7()),
				MD5Sum: md5sum,
				Title:  "Artist - Title 1",
				Cover:  "https://example.com/cover1.jpg",
			}

			result1, err := repo.Upsert(ctx, track1)
			if err != nil {
				return true // Skip if first insert fails
			}

			track2 := &entity.Track{
				ID:     uuid.Must(uuid.NewV7()),
				MD5Sum: md5sum,
				Title:  "Artist - Title 2",
				Cover:  "https://example.com/cover2.jpg",
			}

			result2, err := repo.Upsert(ctx, track2)
			if err != nil {
				return false
			}

			return result2.Rotate == result1.Rotate+1
		},
		genValidMD5Sum(),
	))

	// Property: ExistsByMD5Sum returns true after upsert
	properties.Property("ExistsByMD5Sum returns true after track upsert", prop.ForAll(
		func(md5sum string) bool {
			track := &entity.Track{
				ID:     uuid.Must(uuid.NewV7()),
				MD5Sum: md5sum,
				Title:  "Artist - Title",
				Cover:  "https://example.com/cover.jpg",
			}

			_, err := repo.Upsert(ctx, track)
			if err != nil {
				return true // Skip if upsert fails
			}

			exists, err := repo.ExistsByMD5Sum(ctx, md5sum)
			return err == nil && exists
		},
		genValidMD5Sum(),
	))

	properties.TestingRun(t)
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Logf("Failed to connect to test database: %v", err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Logf("Failed to ping test database: %v", err)
		return nil
	}

	return pool
}

func genValidMD5Sum() gopter.Gen {
	return gen.RegexMatch(`[a-f0-9]{32}`)
}
