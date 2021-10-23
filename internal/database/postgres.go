package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(connString string) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) GetDbPool() *pgxpool.Pool {
	return p.pool
}

func (p *Postgres) GetConnection() (*pgxpool.Conn, error) {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
