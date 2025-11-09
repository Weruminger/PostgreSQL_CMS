//go:build integration

package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func applySQLDir(t *testing.T, pool *pgxpool.Pool, dir string) {
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	var files []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	ctx := context.Background()
	for _, f := range files {
		p := filepath.Join(dir, f)
		b, err := os.ReadFile(p)
		require.NoError(t, err, "read %s", p)
		_, err = pool.Exec(ctx, string(b))
		require.NoError(t, err, "exec %s", p)
	}
}

func TestDB_RLS_Seed_EndToEnd(t *testing.T) {
	ctx := context.Background()

	pgC, err := postgres.RunContainer(ctx,
		postgres.WithInitScripts(), // we apply manually
		postgres.WithDatabase("cms"),
		postgres.WithUsername("cms_owner"),
		postgres.WithPassword("changeme"),
		postgres.WithWaitStrategy(tc.Waiter{Timeout: 60 * time.Second}),
	)
	require.NoError(t, err)
	defer pgC.Terminate(ctx)

	host, _ := pgC.Host(ctx)
	port, _ := pgC.MappedPort(ctx, "5432/tcp")
	url := fmt.Sprintf("postgres://cms_owner:changeme@%s:%s/cms", host, port.Port())

	cfg, err := pgxpool.ParseConfig(url)
	require.NoError(t, err)
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	defer pool.Close()

	// create app role
	_, err = pool.Exec(ctx, `do $$ begin if not exists (select from pg_roles where rolname='cms_app') then create role cms_app login password 'changeme'; end if; end $$;`)
	require.NoError(t, err)

	// apply migrations
	applySQLDir(t, pool, filepath.Join("..", "db"))

	// sanity: seed entry is visible via view
	var slug string
	err = pool.QueryRow(ctx, `select slug from api_published_articles limit 1`).Scan(&slug)
	require.NoError(t, err)
	require.NotEmpty(t, slug)
}
