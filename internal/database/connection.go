package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The Database interface abstracts our database layer, allowing us to swap
// implementations without changing our core application logic.
type Database interface {
	GetConnection(ctx context.Context, key string) (*pgxpool.Pool, error)
	Close()
}

// localDB represents a connection to a single database.
type localDB struct {
	pool *pgxpool.Pool
}

func (ld *localDB) GetConnection(ctx context.Context, key string) (*pgxpool.Pool, error) {

	return ld.pool, nil
}

func (ld *localDB) Close() {
	if ld.pool != nil {
		ld.pool.Close()
		fmt.Println("Closed local database connection.")
	}
}

// shardedDB represents a pool of connections to multiple database shards.
type shardedDB struct {
	shards []*pgxpool.Pool
}

func (sd *shardedDB) GetConnection(ctx context.Context, key string) (*pgxpool.Pool, error) {
	// In a real application, this function would contain sharding logic
	// to pick the correct shard based on the key.
	// For this example, we'll just pick the first shard.
	hash := sha256.Sum256([]byte(key))
	index := int(hash[0]) % len(sd.shards)
	if key == "index" {
		return sd.shards[len(sd.shards)-1], nil // Return the last shard for index queries
	}
	if len(sd.shards) == 0 {
		return nil, fmt.Errorf("no database shards configured")
	}
	return sd.shards[index], nil
}

func (sd *shardedDB) Close() {
	for _, pool := range sd.shards {
		if pool != nil {
			pool.Close()
		}
	}
	fmt.Println("Closed sharded database connections.")
}

// initDatabase is the core of the configuration logic.
// It checks the ENV variable and returns the correct Database implementation.
func initDatabase() (Database, error) {

	env := os.Getenv("ENV")

	if env == "local" || env == "" {
		// Local mode: use a single connection string
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			fmt.Println("DATABASE_URL not set, using default for local development...")
			connStr = "postgres://postgres:admin@localhost:5432/URL?sslmode=disable"
		}
		pool, err := pgxpool.New(context.Background(), connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to local database: %w", err)
		}
		return &localDB{pool: pool}, nil
	}

	// Cloud/Production mode: load multiple connection strings
	var shards []*pgxpool.Pool
	for i := 1; ; i++ {
		key := fmt.Sprintf("DATABASE_SHARD_%d_URL", i)
		connStr := os.Getenv(key)
		if connStr == "" {
			break // No more shards to configure
		}
		pool, err := pgxpool.New(context.Background(), connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to shard %d: %w", i, err)
		}
		shards = append(shards, pool)
	}
	connStr := os.Getenv("DATABASE_INDEX_URL")
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to index database: %w", err)
	} else {
		shards = append(shards, pool)
	}
	if len(shards) == 0 {
		return nil, fmt.Errorf("no database shards configured for cloud environment")
	}

	fmt.Printf("Successfully connected to %d database shards.", len(shards))
	return &shardedDB{shards: shards}, nil
}
