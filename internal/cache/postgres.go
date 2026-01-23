package cache

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	defaultTable      = "cache_kv"
	defaultTTLMinutes = 5
)

var identifierPattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// PostgresStore implements Store using an unlogged hstore table.
type PostgresStore struct {
	conn       *pgx.Conn
	table      string
	defaultTTL time.Duration
}

// NewPostgresStore constructs a Postgres-backed cache store.
func NewPostgresStore(conn *pgx.Conn, opts Options) *PostgresStore {
	return &PostgresStore{
		conn:       conn,
		table:      sanitizeIdentifier(opts.Table),
		defaultTTL: resolveDefaultTTL(opts.DefaultTTL),
	}
}

// EnsureSchema creates the required extension and unlogged cache table.
func (s *PostgresStore) EnsureSchema(ctx context.Context) error {
	if s.conn == nil {
		return errors.New("postgres cache store requires a connection")
	}
	table := s.table
	if table == "" {
		table = defaultTable
	}
	statements := []string{
		"CREATE EXTENSION IF NOT EXISTS hstore",
		fmt.Sprintf(`CREATE UNLOGGED TABLE IF NOT EXISTS %s (
  key text PRIMARY KEY,
  value hstore NOT NULL,
  expires_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
)`, table),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s_expires_at_idx ON %s (expires_at)", table, table),
	}
	for _, statement := range statements {
		if _, err := s.conn.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

// Get returns the cached value for the key when present and not expired.
func (s *PostgresStore) Get(ctx context.Context, key string) (string, bool, error) {
	if s.conn == nil {
		return "", false, errors.New("postgres cache store requires a connection")
	}
	var (
		payload   pgtype.Hstore
		expiresAt time.Time
	)
	query := fmt.Sprintf("SELECT value, expires_at FROM %s WHERE key = $1", s.tableName())
	row := s.conn.QueryRow(ctx, query, key)
	if err := row.Scan(&payload, &expiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	if time.Now().After(expiresAt) {
		_ = s.Delete(ctx, key)
		return "", false, nil
	}
	if payload == nil {
		return "", false, nil
	}
	value, ok := payload["value"]
	if !ok || value == nil {
		return "", false, nil
	}
	return *value, true, nil
}

// Set writes the value to the cache with the provided TTL.
func (s *PostgresStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if s.conn == nil {
		return errors.New("postgres cache store requires a connection")
	}
	if ttl <= 0 {
		ttl = s.defaultTTL
	}
	expiresAt := time.Now().Add(ttl)
	query := fmt.Sprintf(`INSERT INTO %s (key, value, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, expires_at = EXCLUDED.expires_at, updated_at = now()`, s.tableName())
	payload := pgtype.Hstore{"value": &value}
	_, err := s.conn.Exec(ctx, query, key, payload, expiresAt)
	return err
}

// Delete removes a cached key.
func (s *PostgresStore) Delete(ctx context.Context, key string) error {
	if s.conn == nil {
		return errors.New("postgres cache store requires a connection")
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE key = $1", s.tableName())
	_, err := s.conn.Exec(ctx, query, key)
	return err
}

// PurgeExpired removes expired cache entries.
func (s *PostgresStore) PurgeExpired(ctx context.Context) (int64, error) {
	if s.conn == nil {
		return 0, errors.New("postgres cache store requires a connection")
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE expires_at <= now()", s.tableName())
	res, err := s.conn.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (s *PostgresStore) tableName() string {
	if s.table == "" {
		return defaultTable
	}
	return s.table
}

func resolveDefaultTTL(ttl time.Duration) time.Duration {
	if ttl > 0 {
		return ttl
	}
	return time.Duration(defaultTTLMinutes) * time.Minute
}

func sanitizeIdentifier(name string) string {
	if name == "" {
		return ""
	}
	if !identifierPattern.MatchString(name) {
		return ""
	}
	return name
}
