package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Database interface {
	Reader
	Writer
	Migration
}

type Reader interface {
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Select(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Writer interface {
	Insert(ctx context.Context, query string, args ...interface{}) (int64, bool)
	Update(ctx context.Context, query string, args ...interface{}) bool
	Delete(ctx context.Context, query string, args ...interface{}) bool
}

type Migration interface {
	RunSchema(schemas []string, logger *logrus.Logger)
	RebuildIndexes(logger *logrus.Logger)
	Insert(ctx context.Context, query string, args ...interface{}) (int64, bool)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Select(ctx context.Context, query string, args ...interface{}) *sql.Row
	Update(ctx context.Context, query string, args ...interface{}) bool
	Delete(ctx context.Context, query string, args ...interface{}) bool
}

type DatabaseParam struct {
	ConnectionStr string
	DatabaseName  string
	MaxLifetime   time.Duration
	MaxIdleConns  int
	MaxOpenConns  int
}

// postgresql struct
type postgresql struct {
	database *sql.DB
	param    *DatabaseParam
}

func NewConnection(param DatabaseParam, logger *logrus.Logger) Database {
	return &postgresql{
		database: connect(param, logger),
		param:    &param,
	}
}

// connect creates a database connection
func connect(param DatabaseParam, logger *logrus.Logger) *sql.DB {
	logger.Info("Creating database connection")
	conn, err := sql.Open("postgres", param.ConnectionStr)
	if err != nil {
		logger.Fatal(err)
	}

	// confirm connection
	err = conn.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	conn.SetConnMaxLifetime(param.MaxLifetime)
	conn.SetMaxIdleConns(param.MaxIdleConns)
	conn.SetMaxOpenConns(param.MaxOpenConns)
	logger.Info("Database connection successful")
	return conn
}

// RunSchema prepare and execute database changes
func (p *postgresql) RunSchema(schemas []string, logger *logrus.Logger) {
	logger.Info("Executing database schema")
	for _, schema := range schemas {
		if _, err := p.database.Exec(schema); err != nil {
			logger.Error(fmt.Sprintf("Executing database schema failed, because of %v", err))
		}
	}
	logger.Info("Executing database schema completed")
}

// RebuildIndexes execute indexes
func (p *postgresql) RebuildIndexes(logger *logrus.Logger) {
	logger.Info("Rebuild database indexes")
	script := string("REINDEX DATABASE " + p.param.DatabaseName + ";")
	if _, err := p.database.Exec(script); err != nil {
		logger.Error(fmt.Sprintf("Rebuild database indexes failed, because of %v", err))
	}
	logger.Info("Rebuild database indexes completed")
}

// Insert method make a single row query to the databases
func (p *postgresql) Insert(ctx context.Context, query string, args ...interface{}) (int64, bool) {
	stmt, err := p.database.PrepareContext(ctx, query)
	if err != nil {
		return 0, false
	}
	defer func() {
		_ = stmt.Close()
	}()
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, false
	}
	rows, err := res.RowsAffected()
	if err != nil || rows < 1 {
		return 0, false
	}
	return rows, true
}

// Query method make a resultset rows query to the databases
func (p *postgresql) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := p.database.PrepareContext(ctx, query)
	if err != nil {
		return &sql.Rows{}, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	return stmt.QueryContext(ctx, args...)
}

// Select method make a single row query to the databases
func (p *postgresql) Select(ctx context.Context, query string, args ...interface{}) *sql.Row {
	stmt, err := p.database.PrepareContext(ctx, query)
	if err != nil {
		return &sql.Row{}
	}
	defer func() {
		_ = stmt.Close()
	}()
	rows := stmt.QueryRowContext(ctx, args...)
	return rows
}

// Update method executes update database changes to the databases
func (p *postgresql) Update(ctx context.Context, query string, args ...interface{}) bool {
	return execute(p, query, ctx, args)
}

// Delete method executes delete database changes to the databases
func (p *postgresql) Delete(ctx context.Context, query string, args ...interface{}) bool {
	return execute(p, query, ctx, args)
}

// update or delete records
func execute(p *postgresql, query string, ctx context.Context, args []interface{}) bool {
	stmt, err := p.database.PrepareContext(ctx, query)
	if err != nil {
		return false
	}
	defer func() {
		_ = stmt.Close()
	}()
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return false
	}

	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		return false
	}
	return true
}
