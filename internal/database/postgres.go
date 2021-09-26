package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(connString string) (*Postgres, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
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

func (p *Postgres) GetConnection() *pgxpool.Conn {
	conn, err := p.pool.Acquire(context.Background())
	if err != nil {
		fmt.Println("Unable to acquire a database connection", err.Error())
		return nil
	}

	return conn
}
