package internal
import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)
type DB struct { pool *pgxpool.Pool }
func NewDBFromEnv() (*DB, error) {
	u := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGDATABASE"))
	cfg, err := pgxpool.ParseConfig(u); if err != nil { return nil, err }
	cfg.MaxConns = 10; cfg.HealthCheckPeriod = 30 * time.Second
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg); if err != nil { return nil, err }
	return &DB{pool}, nil
}
func (d *DB) Close(){ d.pool.Close() }
func (d *DB) Pool() *pgxpool.Pool { return d.pool }
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) interface{ Scan(dest ...any) error }
	Query(ctx context.Context, sql string, args ...any) (interface{ Next() bool; Scan(dest ...any) error; Err() error; Close() }, error)
	Exec(ctx context.Context, sql string, args ...any) (interface{ String() string }, error)
}
func (d *DB) WithAppSettings(ctx context.Context, tenantID int64, role string, fn func(ctx context.Context, q Querier) error) error {
	conn, err := d.pool.Acquire(ctx); if err != nil { return err }
	defer conn.Release()
	_, err = conn.Exec(ctx, "set local app.tenant_id=$1; set local app.role=$2;", tenantID, role); if err != nil { return err }
	return fn(ctx, conn)
}
