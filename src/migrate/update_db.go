package migrate

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

// Update the database
func Update() {
	connStr := os.Getenv("DB_CONN_STR")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	defer func(ctx2 context.Context, closeCon *pgx.Conn) {
		err := conn.Close(ctx2)
		if err != nil {
			log.Printf("Close DB conneciton failed: %v", err)
		}
	}(context.Background(), conn)

	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	var now time.Time
	tag, err := conn.Exec(ctx, "CREATE DATABASE rotavator;")
	fmt.Printf("rows %v", tag.RowsAffected())
	//err = conn.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatal("failed to execute query", err)
	}

	fmt.Println(now)
}
