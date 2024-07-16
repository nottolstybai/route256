package db

import "github.com/jackc/pgx/v5/pgxpool"

type ConnectionPool struct {
	master  *pgxpool.Pool
	replica *pgxpool.Pool
	current int
}

func NewConnectionPool(master *pgxpool.Pool, replica *pgxpool.Pool) *ConnectionPool {
	return &ConnectionPool{
		master:  master,
		replica: replica,
		current: 0,
	}
}

// Acquire returns connection to master or replica using round-robin
func (p *ConnectionPool) Acquire() *pgxpool.Pool {
	p.current++
	if p.current%2 == 0 {
		return p.master
	}
	return p.replica
}

// Master returns a connection to master
func (p *ConnectionPool) Master() *pgxpool.Pool {
	return p.master
}

func (p *ConnectionPool) Close() {
	p.master.Close()
	p.replica.Close()
}
