package main

import (
	"context"
	"database/sql"
	"fmt"

	pool "github.com/harshadixit12/monolith/src/connection-pool"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "testdb"
)

func query(connPool *pool.ConnectionPool) {
	conn := connPool.Get(context.TODO())
	defer connPool.Put(context.TODO(), conn)

	rows := conn.QueryRowContext(context.TODO(), "SELECT year FROM cars;")

	fmt.Println(rows)
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT year FROM cars;")
	fmt.Println(rows)

	connpool, err := pool.NewConnectionPool(context.TODO(), pool.ConnectionPoolConfig{Size: 10, Timeout: 30000, DB: db})
	if err != nil {
		panic(err)
	}

	//fmt.Println(connpool)

	for i := 0; i < 10; i++ {
		go query(connpool)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
