// Code generated by pggen. DO NOT EDIT.

package pg

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
	FindEnumTypes(ctx context.Context, oIDs []uint32) ([]FindEnumTypesRow, error)
	// FindEnumTypesBatch enqueues a FindEnumTypes query into batch to be executed
	// later by the batch.
	FindEnumTypesBatch(ctx context.Context, batch *pgx.Batch, oIDs []uint32)
	// FindEnumTypesScan scans the result of an executed FindEnumTypesBatch query.
	FindEnumTypesScan(results pgx.BatchResults) ([]FindEnumTypesRow, error)
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

const findEnumTypesSQL = `SELECT typ.oid            AS oid,
       typ.typname        AS type_name,
       typ.typtype        AS type_kind,
       enum.oid           AS enum_oid,
       enum.enumsortorder AS enum_order,
       enum.enumlabel     AS enum_label,
       typ.typelem        AS elem_type,
       typ.typarray       AS array_type,
       typ.typrelid       AS composite_type_id,
       typ.typnotnull     AS domain_not_null_constraint,
       typ.typbasetype    AS domain_base_type,
       typ.typndims       AS num_dimensions,
       typ.typdefault     AS default_expr
FROM pg_type typ
  JOIN pg_enum enum ON typ.oid = enum.enumtypid
WHERE typ.typisdefined
  AND typ.typtype = 'e'AND typ.oid = ANY ($1::oid[])
ORDER BY typ.oid DESC;`

type FindEnumTypesRow struct {
	OID                     pgtype.OID    `json:"oid"`
	TypeName                pgtype.Name   `json:"type_name"`
	TypeKind                pgtype.QChar  `json:"type_kind"`
	EnumOID                 pgtype.OID    `json:"enum_oid"`
	EnumOrder               pgtype.Float4 `json:"enum_order"`
	EnumLabel               pgtype.Name   `json:"enum_label"`
	ElemType                pgtype.OID    `json:"elem_type"`
	ArrayType               pgtype.OID    `json:"array_type"`
	CompositeTypeID         pgtype.OID    `json:"composite_type_id"`
	DomainNotNullConstraint pgtype.Bool   `json:"domain_not_null_constraint"`
	DomainBaseType          pgtype.OID    `json:"domain_base_type"`
	NumDimensions           pgtype.Int4   `json:"num_dimensions"`
	DefaultExpr             pgtype.Text   `json:"default_expr"`
}

// FindEnumTypes implements Querier.FindEnumTypes.
func (q *DBQuerier) FindEnumTypes(ctx context.Context, oIDs []uint32) ([]FindEnumTypesRow, error) {
	rows, err := q.conn.Query(ctx, findEnumTypesSQL, oIDs)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query FindEnumTypes: %w", err)
	}
	items := []FindEnumTypesRow{}
	for rows.Next() {
		var item FindEnumTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.TypeKind, &item.EnumOID, &item.EnumOrder, &item.EnumLabel, &item.ElemType, &item.ArrayType, &item.CompositeTypeID, &item.DomainNotNullConstraint, &item.DomainBaseType, &item.NumDimensions, &item.DefaultExpr); err != nil {
			return nil, fmt.Errorf("scan FindEnumTypes row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// FindEnumTypesBatch implements Querier.FindEnumTypesBatch.
func (q *DBQuerier) FindEnumTypesBatch(ctx context.Context, batch *pgx.Batch, oIDs []uint32) {
	batch.Queue(findEnumTypesSQL, oIDs)
}

// FindEnumTypesScan implements Querier.FindEnumTypesScan.
func (q *DBQuerier) FindEnumTypesScan(results pgx.BatchResults) ([]FindEnumTypesRow, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	items := []FindEnumTypesRow{}
	for rows.Next() {
		var item FindEnumTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.TypeKind, &item.EnumOID, &item.EnumOrder, &item.EnumLabel, &item.ElemType, &item.ArrayType, &item.CompositeTypeID, &item.DomainNotNullConstraint, &item.DomainBaseType, &item.NumDimensions, &item.DefaultExpr); err != nil {
			return nil, fmt.Errorf("scan FindEnumTypesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}