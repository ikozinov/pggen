// Code generated by pggen. DO NOT EDIT.

package device

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	FindDevicesByUser(ctx context.Context, id int) ([]FindDevicesByUserRow, error)
	// FindDevicesByUserBatch enqueues a FindDevicesByUser query into batch to be executed
	// later by the batch.
	FindDevicesByUserBatch(batch *pgx.Batch, id int)
	// FindDevicesByUserScan scans the result of an executed FindDevicesByUserBatch query.
	FindDevicesByUserScan(results pgx.BatchResults) ([]FindDevicesByUserRow, error)

	CompositeUser(ctx context.Context) ([]CompositeUserRow, error)
	// CompositeUserBatch enqueues a CompositeUser query into batch to be executed
	// later by the batch.
	CompositeUserBatch(batch *pgx.Batch)
	// CompositeUserScan scans the result of an executed CompositeUserBatch query.
	CompositeUserScan(results pgx.BatchResults) ([]CompositeUserRow, error)

	CompositeUserOne(ctx context.Context) (User, error)
	// CompositeUserOneBatch enqueues a CompositeUserOne query into batch to be executed
	// later by the batch.
	CompositeUserOneBatch(batch *pgx.Batch)
	// CompositeUserOneScan scans the result of an executed CompositeUserOneBatch query.
	CompositeUserOneScan(results pgx.BatchResults) (User, error)

	CompositeUserOneTwoCols(ctx context.Context) (CompositeUserOneTwoColsRow, error)
	// CompositeUserOneTwoColsBatch enqueues a CompositeUserOneTwoCols query into batch to be executed
	// later by the batch.
	CompositeUserOneTwoColsBatch(batch *pgx.Batch)
	// CompositeUserOneTwoColsScan scans the result of an executed CompositeUserOneTwoColsBatch query.
	CompositeUserOneTwoColsScan(results pgx.BatchResults) (CompositeUserOneTwoColsRow, error)

	CompositeUserMany(ctx context.Context) ([]User, error)
	// CompositeUserManyBatch enqueues a CompositeUserMany query into batch to be executed
	// later by the batch.
	CompositeUserManyBatch(batch *pgx.Batch)
	// CompositeUserManyScan scans the result of an executed CompositeUserManyBatch query.
	CompositeUserManyScan(results pgx.BatchResults) ([]User, error)

	InsertUser(ctx context.Context, userID int, name string) (pgconn.CommandTag, error)
	// InsertUserBatch enqueues a InsertUser query into batch to be executed
	// later by the batch.
	InsertUserBatch(batch *pgx.Batch, userID int, name string)
	// InsertUserScan scans the result of an executed InsertUserBatch query.
	InsertUserScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	InsertDevice(ctx context.Context, mac pgtype.Macaddr, owner int) (pgconn.CommandTag, error)
	// InsertDeviceBatch enqueues a InsertDevice query into batch to be executed
	// later by the batch.
	InsertDeviceBatch(batch *pgx.Batch, mac pgtype.Macaddr, owner int)
	// InsertDeviceScan scans the result of an executed InsertDeviceBatch query.
	InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error)
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

// User represents the Postgres composite type "user".
type User struct {
	ID   pgtype.Int8 `json:"id"`
	Name pgtype.Text `json:"name"`
}

// DeviceType represents the Postgres enum "device_type".
type DeviceType string

const (
	DeviceTypeUndefined DeviceType = "undefined"
	DeviceTypePhone     DeviceType = "phone"
	DeviceTypeLaptop    DeviceType = "laptop"
	DeviceTypeIpad      DeviceType = "ipad"
	DeviceTypeDesktop   DeviceType = "desktop"
	DeviceTypeIot       DeviceType = "iot"
)

func (d DeviceType) String() string { return string(d) }

const findDevicesByUserSQL = `SELECT
  id,
  name,
  (SELECT array_agg(mac) FROM device WHERE owner = id)
FROM "user"
WHERE id = $1;`

type FindDevicesByUserRow struct {
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	ArrayAgg pgtype.MacaddrArray `json:"array_agg"`
}

// FindDevicesByUser implements Querier.FindDevicesByUser.
func (q *DBQuerier) FindDevicesByUser(ctx context.Context, id int) ([]FindDevicesByUserRow, error) {
	rows, err := q.conn.Query(ctx, findDevicesByUserSQL, id)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query FindDevicesByUser: %w", err)
	}
	items := []FindDevicesByUserRow{}
	for rows.Next() {
		var item FindDevicesByUserRow
		if err := rows.Scan(&item.ID, &item.Name, &item.ArrayAgg); err != nil {
			return nil, fmt.Errorf("scan FindDevicesByUser row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// FindDevicesByUserBatch implements Querier.FindDevicesByUserBatch.
func (q *DBQuerier) FindDevicesByUserBatch(batch *pgx.Batch, id int) {
	batch.Queue(findDevicesByUserSQL, id)
}

// FindDevicesByUserScan implements Querier.FindDevicesByUserScan.
func (q *DBQuerier) FindDevicesByUserScan(results pgx.BatchResults) ([]FindDevicesByUserRow, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	items := []FindDevicesByUserRow{}
	for rows.Next() {
		var item FindDevicesByUserRow
		if err := rows.Scan(&item.ID, &item.Name, &item.ArrayAgg); err != nil {
			return nil, fmt.Errorf("scan FindDevicesByUserBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

const compositeUserSQL = `SELECT
  d.mac,
  d.type,
  ROW (u.id, u.name)::"user" AS "user"
FROM device d
  LEFT JOIN "user" u ON u.id = d.owner;`

type CompositeUserRow struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
	User User           `json:"user"`
}

// CompositeUser implements Querier.CompositeUser.
func (q *DBQuerier) CompositeUser(ctx context.Context) ([]CompositeUserRow, error) {
	rows, err := q.conn.Query(ctx, compositeUserSQL)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query CompositeUser: %w", err)
	}
	items := []CompositeUserRow{}
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	for rows.Next() {
		var item CompositeUserRow
		if err := rows.Scan(&item.Mac, &item.Type, userRow); err != nil {
			return nil, fmt.Errorf("scan CompositeUser row: %w", err)
		}
		item.User.ID = *userRow[0].(*pgtype.Int8)
		item.User.Name = *userRow[1].(*pgtype.Text)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// CompositeUserBatch implements Querier.CompositeUserBatch.
func (q *DBQuerier) CompositeUserBatch(batch *pgx.Batch) {
	batch.Queue(compositeUserSQL)
}

// CompositeUserScan implements Querier.CompositeUserScan.
func (q *DBQuerier) CompositeUserScan(results pgx.BatchResults) ([]CompositeUserRow, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	items := []CompositeUserRow{}
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	for rows.Next() {
		var item CompositeUserRow
		if err := rows.Scan(&item.Mac, &item.Type, userRow); err != nil {
			return nil, fmt.Errorf("scan CompositeUserBatch row: %w", err)
		}
		item.User.ID = *userRow[0].(*pgtype.Int8)
		item.User.Name = *userRow[1].(*pgtype.Text)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

const compositeUserOneSQL = `SELECT ROW (15, 'qux')::"user" AS "user";`

// CompositeUserOne implements Querier.CompositeUserOne.
func (q *DBQuerier) CompositeUserOne(ctx context.Context) (User, error) {
	row := q.conn.QueryRow(ctx, compositeUserOneSQL)
	var item User
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	if err := row.Scan(userRow); err != nil {
		return item, fmt.Errorf("query CompositeUserOne: %w", err)
	}
	item.ID = *userRow[0].(*pgtype.Int8)
	item.Name = *userRow[1].(*pgtype.Text)
	return item, nil
}

// CompositeUserOneBatch implements Querier.CompositeUserOneBatch.
func (q *DBQuerier) CompositeUserOneBatch(batch *pgx.Batch) {
	batch.Queue(compositeUserOneSQL)
}

// CompositeUserOneScan implements Querier.CompositeUserOneScan.
func (q *DBQuerier) CompositeUserOneScan(results pgx.BatchResults) (User, error) {
	row := results.QueryRow()
	var item User
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	if err := row.Scan(userRow); err != nil {
		return item, fmt.Errorf("scan CompositeUserOneBatch row: %w", err)
	}
	item.ID = *userRow[0].(*pgtype.Int8)
	item.Name = *userRow[1].(*pgtype.Text)
	return item, nil
}

const compositeUserOneTwoColsSQL = `SELECT 1 AS num, ROW (15, 'qux')::"user" AS "user";`

type CompositeUserOneTwoColsRow struct {
	Num  int32 `json:"num"`
	User User  `json:"user"`
}

// CompositeUserOneTwoCols implements Querier.CompositeUserOneTwoCols.
func (q *DBQuerier) CompositeUserOneTwoCols(ctx context.Context) (CompositeUserOneTwoColsRow, error) {
	row := q.conn.QueryRow(ctx, compositeUserOneTwoColsSQL)
	var item CompositeUserOneTwoColsRow
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	if err := row.Scan(&item.Num, userRow); err != nil {
		return item, fmt.Errorf("query CompositeUserOneTwoCols: %w", err)
	}
	item.User.ID = *userRow[0].(*pgtype.Int8)
	item.User.Name = *userRow[1].(*pgtype.Text)
	return item, nil
}

// CompositeUserOneTwoColsBatch implements Querier.CompositeUserOneTwoColsBatch.
func (q *DBQuerier) CompositeUserOneTwoColsBatch(batch *pgx.Batch) {
	batch.Queue(compositeUserOneTwoColsSQL)
}

// CompositeUserOneTwoColsScan implements Querier.CompositeUserOneTwoColsScan.
func (q *DBQuerier) CompositeUserOneTwoColsScan(results pgx.BatchResults) (CompositeUserOneTwoColsRow, error) {
	row := results.QueryRow()
	var item CompositeUserOneTwoColsRow
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	if err := row.Scan(&item.Num, userRow); err != nil {
		return item, fmt.Errorf("scan CompositeUserOneTwoColsBatch row: %w", err)
	}
	item.User.ID = *userRow[0].(*pgtype.Int8)
	item.User.Name = *userRow[1].(*pgtype.Text)
	return item, nil
}

const compositeUserManySQL = `SELECT ROW (15, 'qux')::"user" AS "user";`

// CompositeUserMany implements Querier.CompositeUserMany.
func (q *DBQuerier) CompositeUserMany(ctx context.Context) ([]User, error) {
	rows, err := q.conn.Query(ctx, compositeUserManySQL)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query CompositeUserMany: %w", err)
	}
	items := []User{}
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	for rows.Next() {
		var item User
		if err := rows.Scan(userRow); err != nil {
			return nil, fmt.Errorf("scan CompositeUserMany row: %w", err)
		}
		item.ID = *userRow[0].(*pgtype.Int8)
		item.Name = *userRow[1].(*pgtype.Text)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// CompositeUserManyBatch implements Querier.CompositeUserManyBatch.
func (q *DBQuerier) CompositeUserManyBatch(batch *pgx.Batch) {
	batch.Queue(compositeUserManySQL)
}

// CompositeUserManyScan implements Querier.CompositeUserManyScan.
func (q *DBQuerier) CompositeUserManyScan(results pgx.BatchResults) ([]User, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	items := []User{}
	userRow := pgtype.CompositeFields{
		&pgtype.Int8{},
		&pgtype.Text{},
	}
	for rows.Next() {
		var item User
		if err := rows.Scan(userRow); err != nil {
			return nil, fmt.Errorf("scan CompositeUserManyBatch row: %w", err)
		}
		item.ID = *userRow[0].(*pgtype.Int8)
		item.Name = *userRow[1].(*pgtype.Text)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

const insertUserSQL = `INSERT INTO "user" (id, name) VALUES ($1, $2);`

// InsertUser implements Querier.InsertUser.
func (q *DBQuerier) InsertUser(ctx context.Context, userID int, name string) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, insertUserSQL, userID, name)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertUser: %w", err)
	}
	return cmdTag, err
}

// InsertUserBatch implements Querier.InsertUserBatch.
func (q *DBQuerier) InsertUserBatch(batch *pgx.Batch, userID int, name string) {
	batch.Queue(insertUserSQL, userID, name)
}

// InsertUserScan implements Querier.InsertUserScan.
func (q *DBQuerier) InsertUserScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertUserBatch: %w", err)
	}
	return cmdTag, err
}

const insertDeviceSQL = `INSERT INTO device (mac, owner) VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac pgtype.Macaddr, owner int) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, owner)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
}

// InsertDeviceBatch implements Querier.InsertDeviceBatch.
func (q *DBQuerier) InsertDeviceBatch(batch *pgx.Batch, mac pgtype.Macaddr, owner int) {
	batch.Queue(insertDeviceSQL, mac, owner)
}

// InsertDeviceScan implements Querier.InsertDeviceScan.
func (q *DBQuerier) InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertDeviceBatch: %w", err)
	}
	return cmdTag, err
}
