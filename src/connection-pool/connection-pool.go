package pool

import (
	"context"
	"database/sql"

	customds "github.com/harshadixit12/monolith/shared/custom-ds"
	_ "github.com/lib/pq"
)

type ConnectionPoolConfig struct {
	Size    int
	Timeout int
	DB      *sql.DB
}

type ConnectionPool struct {
	size    int
	pool    customds.BlockingQueue[sql.Conn]
	db      *sql.DB
	timeout int
}

func NewConnectionPool(ctx context.Context, config ConnectionPoolConfig) (*ConnectionPool, error) {
	connectionPool := &ConnectionPool{
		size: config.Size, db: config.DB,
		pool: *customds.NewBlockingBlockingQueue[sql.Conn](config.Size), timeout: config.Timeout}

	for i := 0; i < config.Size; i++ {
		sqlConn, err := config.DB.Conn(ctx)

		if err != nil {
			return nil, err
		}
		connectionPool.pool.Put(*sqlConn)
	}

	return connectionPool, nil
}

func (p *ConnectionPool) Get(ctx context.Context) *sql.Conn {
	conn := p.pool.Take()

	return &conn
}

func (p *ConnectionPool) Put(ctx context.Context, conn *sql.Conn) {
	p.pool.Put(*conn)
}
