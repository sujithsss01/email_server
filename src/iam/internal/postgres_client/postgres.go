package postgresclient

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sujithsss01/email_server/config"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDatabaseUrl(config *config.Config) string {
	url := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.PSQLConfig.Username, config.PSQLConfig.Password),
		Host:   net.JoinHostPort(config.PSQLConfig.Hostname, config.PSQLConfig.Port),
		Path:   config.PSQLConfig.DBNAME,
	}
	query := url.Query()
	query.Set("sslmode", "disable")
	url.RawQuery = query.Encode()
	databaseURL := url.String()
	return databaseURL
}

func NewConnection(ctx context.Context) *sql.DB {
	cf := config.GetConfig()
	databaseUrl := getDatabaseUrl(cf)

	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		panic(err)
	}
	return db
}

func NewGORM(sqlDB *sql.DB) *gorm.DB {
	db, err := gorm.Open(
		gormpostgres.New(gormpostgres.Config{
			Conn: sqlDB,
		}),
	)
	if err != nil {
		panic(err)
	}
	return db
}

func RunMigrations() {
	cf := config.GetConfig()
	sqlDB, err := sql.Open("pgx", getDatabaseUrl(cf))
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	driver, err := migratepostgres.WithInstance(
		sqlDB,
		&migratepostgres.Config{},
	)
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		_ = driver.Close()
		panic(err)
	}

	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("Close migration source")
			panic(sourceErr)
		}
		if databaseErr != nil {
			log.Printf("Close migration database")
			panic(databaseErr)
		}
	}()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	log.Println("IAM migrations completed")
}
