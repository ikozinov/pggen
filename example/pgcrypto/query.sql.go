// Code generated by pggen. DO NOT EDIT.

package pgcrypto

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	CreateUser(ctx context.Context, email string, password string) (pgconn.CommandTag, error)
	// CreateUserBatch enqueues a CreateUser query into batch to be executed
	// later by the batch.
	CreateUserBatch(batch *pgx.Batch, email string, password string)
	// CreateUserScan scans the result of an executed CreateUserBatch query.
	CreateUserScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	FindUser(ctx context.Context, email string) (FindUserRow, error)
	// FindUserBatch enqueues a FindUser query into batch to be executed
	// later by the batch.
	FindUserBatch(batch *pgx.Batch, email string)
	// FindUserScan scans the result of an executed FindUserBatch query.
	FindUserScan(results pgx.BatchResults) (FindUserRow, error)
}

type DBQuerier struct {
	conn genericConn
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
	}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

const createUserSQL = `INSERT INTO "user" (email, pass)
VALUES ($1, crypt($2, gen_salt('bf')));`

// CreateUser implements Querier.CreateUser.
func (q *DBQuerier) CreateUser(ctx context.Context, email string, password string) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, createUserSQL, email, password)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query CreateUser: %w", err)
	}
	return cmdTag, err
}

// CreateUserBatch implements Querier.CreateUserBatch.
func (q *DBQuerier) CreateUserBatch(batch *pgx.Batch, email string, password string) {
	batch.Queue(createUserSQL, email, password)
}

// CreateUserScan implements Querier.CreateUserScan.
func (q *DBQuerier) CreateUserScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec CreateUserBatch: %w", err)
	}
	return cmdTag, err
}

const findUserSQL = `SELECT email, pass from "user"
where email = $1;`

type FindUserRow struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

// FindUser implements Querier.FindUser.
func (q *DBQuerier) FindUser(ctx context.Context, email string) (FindUserRow, error) {
	row := q.conn.QueryRow(ctx, findUserSQL, email)
	var item FindUserRow
	if err := row.Scan(&item.Email, &item.Pass); err != nil {
		return item, fmt.Errorf("query FindUser: %w", err)
	}
	return item, nil
}

// FindUserBatch implements Querier.FindUserBatch.
func (q *DBQuerier) FindUserBatch(batch *pgx.Batch, email string) {
	batch.Queue(findUserSQL, email)
}

// FindUserScan implements Querier.FindUserScan.
func (q *DBQuerier) FindUserScan(results pgx.BatchResults) (FindUserRow, error) {
	row := results.QueryRow()
	var item FindUserRow
	if err := row.Scan(&item.Email, &item.Pass); err != nil {
		return item, fmt.Errorf("scan FindUserBatch row: %w", err)
	}
	return item, nil
}
