package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/emarifer/gocms/settings"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

// NewMariaDBConnection creates a new MariaDB database connection
func NewMariaDBConnection(
	ctx context.Context, s *settings.Settings,
) (*sqlx.DB, error) {
	// connectionString := "root:my-secret-pw@tcp(localhost:3306)/cms_db?parseTime=true" // harcoded connection string

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		s.DB.User,
		s.DB.Password,
		s.DB.Host,
		s.DB.Port,
		s.DB.Name,
	)

	db, err := sqlx.ConnectContext(ctx, "mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	// connection settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// sqlx.BindDriver("mysql", sqlx.DOLLAR)

	log.Println("🚀 Connected Successfully to the Database")

	return db, nil
}
