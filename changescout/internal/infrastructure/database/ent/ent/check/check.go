// Code generated by ent, DO NOT EDIT.

package check

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the check type in the database.
	Label = "check"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldWebsiteID holds the string denoting the website_id field in the database.
	FieldWebsiteID = "website_id"
	// FieldResult holds the string denoting the result field in the database.
	FieldResult = "result"
	// FieldHasError holds the string denoting the has_error field in the database.
	FieldHasError = "has_error"
	// FieldErrorMessage holds the string denoting the error_message field in the database.
	FieldErrorMessage = "error_message"
	// FieldHasDiff holds the string denoting the has_diff field in the database.
	FieldHasDiff = "has_diff"
	// FieldDiffChange holds the string denoting the diff_change field in the database.
	FieldDiffChange = "diff_change"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeWebsite holds the string denoting the website edge name in mutations.
	EdgeWebsite = "website"
	// Table holds the table name of the check in the database.
	Table = "checks"
	// WebsiteTable is the table that holds the website relation/edge.
	WebsiteTable = "checks"
	// WebsiteInverseTable is the table name for the Website entity.
	// It exists in this package in order to avoid circular dependency with the "website" package.
	WebsiteInverseTable = "websites"
	// WebsiteColumn is the table column denoting the website relation/edge.
	WebsiteColumn = "website_id"
)

// Columns holds all SQL columns for check fields.
var Columns = []string{
	FieldID,
	FieldWebsiteID,
	FieldResult,
	FieldHasError,
	FieldErrorMessage,
	FieldHasDiff,
	FieldDiffChange,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// ResultValidator is a validator for the "result" field. It is called by the builders before save.
	ResultValidator func([]byte) error
	// DefaultHasError holds the default value on creation for the "has_error" field.
	DefaultHasError bool
	// DefaultHasDiff holds the default value on creation for the "has_diff" field.
	DefaultHasDiff bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Check queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByWebsiteID orders the results by the website_id field.
func ByWebsiteID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWebsiteID, opts...).ToFunc()
}

// ByHasError orders the results by the has_error field.
func ByHasError(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHasError, opts...).ToFunc()
}

// ByErrorMessage orders the results by the error_message field.
func ByErrorMessage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldErrorMessage, opts...).ToFunc()
}

// ByHasDiff orders the results by the has_diff field.
func ByHasDiff(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHasDiff, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByWebsiteField orders the results by website field.
func ByWebsiteField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newWebsiteStep(), sql.OrderByField(field, opts...))
	}
}
func newWebsiteStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(WebsiteInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, WebsiteTable, WebsiteColumn),
	)
}
