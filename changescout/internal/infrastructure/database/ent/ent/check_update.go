// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/check"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/predicate"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/website"
	"github.com/google/uuid"
)

// CheckUpdate is the builder for updating Check entities.
type CheckUpdate struct {
	config
	hooks    []Hook
	mutation *CheckMutation
}

// Where appends a list predicates to the CheckUpdate builder.
func (cu *CheckUpdate) Where(ps ...predicate.Check) *CheckUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetWebsiteID sets the "website_id" field.
func (cu *CheckUpdate) SetWebsiteID(u uuid.UUID) *CheckUpdate {
	cu.mutation.SetWebsiteID(u)
	return cu
}

// SetNillableWebsiteID sets the "website_id" field if the given value is not nil.
func (cu *CheckUpdate) SetNillableWebsiteID(u *uuid.UUID) *CheckUpdate {
	if u != nil {
		cu.SetWebsiteID(*u)
	}
	return cu
}

// ClearWebsiteID clears the value of the "website_id" field.
func (cu *CheckUpdate) ClearWebsiteID() *CheckUpdate {
	cu.mutation.ClearWebsiteID()
	return cu
}

// SetResult sets the "result" field.
func (cu *CheckUpdate) SetResult(b []byte) *CheckUpdate {
	cu.mutation.SetResult(b)
	return cu
}

// SetHasError sets the "has_error" field.
func (cu *CheckUpdate) SetHasError(b bool) *CheckUpdate {
	cu.mutation.SetHasError(b)
	return cu
}

// SetNillableHasError sets the "has_error" field if the given value is not nil.
func (cu *CheckUpdate) SetNillableHasError(b *bool) *CheckUpdate {
	if b != nil {
		cu.SetHasError(*b)
	}
	return cu
}

// SetErrorMessage sets the "error_message" field.
func (cu *CheckUpdate) SetErrorMessage(s string) *CheckUpdate {
	cu.mutation.SetErrorMessage(s)
	return cu
}

// SetNillableErrorMessage sets the "error_message" field if the given value is not nil.
func (cu *CheckUpdate) SetNillableErrorMessage(s *string) *CheckUpdate {
	if s != nil {
		cu.SetErrorMessage(*s)
	}
	return cu
}

// ClearErrorMessage clears the value of the "error_message" field.
func (cu *CheckUpdate) ClearErrorMessage() *CheckUpdate {
	cu.mutation.ClearErrorMessage()
	return cu
}

// SetHasDiff sets the "has_diff" field.
func (cu *CheckUpdate) SetHasDiff(b bool) *CheckUpdate {
	cu.mutation.SetHasDiff(b)
	return cu
}

// SetNillableHasDiff sets the "has_diff" field if the given value is not nil.
func (cu *CheckUpdate) SetNillableHasDiff(b *bool) *CheckUpdate {
	if b != nil {
		cu.SetHasDiff(*b)
	}
	return cu
}

// SetDiffChange sets the "diff_change" field.
func (cu *CheckUpdate) SetDiffChange(d *diff.Result) *CheckUpdate {
	cu.mutation.SetDiffChange(d)
	return cu
}

// ClearDiffChange clears the value of the "diff_change" field.
func (cu *CheckUpdate) ClearDiffChange() *CheckUpdate {
	cu.mutation.ClearDiffChange()
	return cu
}

// SetCreatedAt sets the "created_at" field.
func (cu *CheckUpdate) SetCreatedAt(t time.Time) *CheckUpdate {
	cu.mutation.SetCreatedAt(t)
	return cu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cu *CheckUpdate) SetNillableCreatedAt(t *time.Time) *CheckUpdate {
	if t != nil {
		cu.SetCreatedAt(*t)
	}
	return cu
}

// SetWebsite sets the "website" edge to the Website entity.
func (cu *CheckUpdate) SetWebsite(w *Website) *CheckUpdate {
	return cu.SetWebsiteID(w.ID)
}

// Mutation returns the CheckMutation object of the builder.
func (cu *CheckUpdate) Mutation() *CheckMutation {
	return cu.mutation
}

// ClearWebsite clears the "website" edge to the Website entity.
func (cu *CheckUpdate) ClearWebsite() *CheckUpdate {
	cu.mutation.ClearWebsite()
	return cu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CheckUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CheckUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CheckUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CheckUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cu *CheckUpdate) check() error {
	if v, ok := cu.mutation.Result(); ok {
		if err := check.ResultValidator(v); err != nil {
			return &ValidationError{Name: "result", err: fmt.Errorf(`ent: validator failed for field "Check.result": %w`, err)}
		}
	}
	return nil
}

func (cu *CheckUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := cu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(check.Table, check.Columns, sqlgraph.NewFieldSpec(check.FieldID, field.TypeUUID))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Result(); ok {
		_spec.SetField(check.FieldResult, field.TypeBytes, value)
	}
	if value, ok := cu.mutation.HasError(); ok {
		_spec.SetField(check.FieldHasError, field.TypeBool, value)
	}
	if value, ok := cu.mutation.ErrorMessage(); ok {
		_spec.SetField(check.FieldErrorMessage, field.TypeString, value)
	}
	if cu.mutation.ErrorMessageCleared() {
		_spec.ClearField(check.FieldErrorMessage, field.TypeString)
	}
	if value, ok := cu.mutation.HasDiff(); ok {
		_spec.SetField(check.FieldHasDiff, field.TypeBool, value)
	}
	if value, ok := cu.mutation.DiffChange(); ok {
		_spec.SetField(check.FieldDiffChange, field.TypeJSON, value)
	}
	if cu.mutation.DiffChangeCleared() {
		_spec.ClearField(check.FieldDiffChange, field.TypeJSON)
	}
	if value, ok := cu.mutation.CreatedAt(); ok {
		_spec.SetField(check.FieldCreatedAt, field.TypeTime, value)
	}
	if cu.mutation.WebsiteCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   check.WebsiteTable,
			Columns: []string{check.WebsiteColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(website.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.WebsiteIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   check.WebsiteTable,
			Columns: []string{check.WebsiteColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(website.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{check.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// CheckUpdateOne is the builder for updating a single Check entity.
type CheckUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CheckMutation
}

// SetWebsiteID sets the "website_id" field.
func (cuo *CheckUpdateOne) SetWebsiteID(u uuid.UUID) *CheckUpdateOne {
	cuo.mutation.SetWebsiteID(u)
	return cuo
}

// SetNillableWebsiteID sets the "website_id" field if the given value is not nil.
func (cuo *CheckUpdateOne) SetNillableWebsiteID(u *uuid.UUID) *CheckUpdateOne {
	if u != nil {
		cuo.SetWebsiteID(*u)
	}
	return cuo
}

// ClearWebsiteID clears the value of the "website_id" field.
func (cuo *CheckUpdateOne) ClearWebsiteID() *CheckUpdateOne {
	cuo.mutation.ClearWebsiteID()
	return cuo
}

// SetResult sets the "result" field.
func (cuo *CheckUpdateOne) SetResult(b []byte) *CheckUpdateOne {
	cuo.mutation.SetResult(b)
	return cuo
}

// SetHasError sets the "has_error" field.
func (cuo *CheckUpdateOne) SetHasError(b bool) *CheckUpdateOne {
	cuo.mutation.SetHasError(b)
	return cuo
}

// SetNillableHasError sets the "has_error" field if the given value is not nil.
func (cuo *CheckUpdateOne) SetNillableHasError(b *bool) *CheckUpdateOne {
	if b != nil {
		cuo.SetHasError(*b)
	}
	return cuo
}

// SetErrorMessage sets the "error_message" field.
func (cuo *CheckUpdateOne) SetErrorMessage(s string) *CheckUpdateOne {
	cuo.mutation.SetErrorMessage(s)
	return cuo
}

// SetNillableErrorMessage sets the "error_message" field if the given value is not nil.
func (cuo *CheckUpdateOne) SetNillableErrorMessage(s *string) *CheckUpdateOne {
	if s != nil {
		cuo.SetErrorMessage(*s)
	}
	return cuo
}

// ClearErrorMessage clears the value of the "error_message" field.
func (cuo *CheckUpdateOne) ClearErrorMessage() *CheckUpdateOne {
	cuo.mutation.ClearErrorMessage()
	return cuo
}

// SetHasDiff sets the "has_diff" field.
func (cuo *CheckUpdateOne) SetHasDiff(b bool) *CheckUpdateOne {
	cuo.mutation.SetHasDiff(b)
	return cuo
}

// SetNillableHasDiff sets the "has_diff" field if the given value is not nil.
func (cuo *CheckUpdateOne) SetNillableHasDiff(b *bool) *CheckUpdateOne {
	if b != nil {
		cuo.SetHasDiff(*b)
	}
	return cuo
}

// SetDiffChange sets the "diff_change" field.
func (cuo *CheckUpdateOne) SetDiffChange(d *diff.Result) *CheckUpdateOne {
	cuo.mutation.SetDiffChange(d)
	return cuo
}

// ClearDiffChange clears the value of the "diff_change" field.
func (cuo *CheckUpdateOne) ClearDiffChange() *CheckUpdateOne {
	cuo.mutation.ClearDiffChange()
	return cuo
}

// SetCreatedAt sets the "created_at" field.
func (cuo *CheckUpdateOne) SetCreatedAt(t time.Time) *CheckUpdateOne {
	cuo.mutation.SetCreatedAt(t)
	return cuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cuo *CheckUpdateOne) SetNillableCreatedAt(t *time.Time) *CheckUpdateOne {
	if t != nil {
		cuo.SetCreatedAt(*t)
	}
	return cuo
}

// SetWebsite sets the "website" edge to the Website entity.
func (cuo *CheckUpdateOne) SetWebsite(w *Website) *CheckUpdateOne {
	return cuo.SetWebsiteID(w.ID)
}

// Mutation returns the CheckMutation object of the builder.
func (cuo *CheckUpdateOne) Mutation() *CheckMutation {
	return cuo.mutation
}

// ClearWebsite clears the "website" edge to the Website entity.
func (cuo *CheckUpdateOne) ClearWebsite() *CheckUpdateOne {
	cuo.mutation.ClearWebsite()
	return cuo
}

// Where appends a list predicates to the CheckUpdate builder.
func (cuo *CheckUpdateOne) Where(ps ...predicate.Check) *CheckUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CheckUpdateOne) Select(field string, fields ...string) *CheckUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Check entity.
func (cuo *CheckUpdateOne) Save(ctx context.Context) (*Check, error) {
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CheckUpdateOne) SaveX(ctx context.Context) *Check {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CheckUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CheckUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cuo *CheckUpdateOne) check() error {
	if v, ok := cuo.mutation.Result(); ok {
		if err := check.ResultValidator(v); err != nil {
			return &ValidationError{Name: "result", err: fmt.Errorf(`ent: validator failed for field "Check.result": %w`, err)}
		}
	}
	return nil
}

func (cuo *CheckUpdateOne) sqlSave(ctx context.Context) (_node *Check, err error) {
	if err := cuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(check.Table, check.Columns, sqlgraph.NewFieldSpec(check.FieldID, field.TypeUUID))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Check.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, check.FieldID)
		for _, f := range fields {
			if !check.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != check.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Result(); ok {
		_spec.SetField(check.FieldResult, field.TypeBytes, value)
	}
	if value, ok := cuo.mutation.HasError(); ok {
		_spec.SetField(check.FieldHasError, field.TypeBool, value)
	}
	if value, ok := cuo.mutation.ErrorMessage(); ok {
		_spec.SetField(check.FieldErrorMessage, field.TypeString, value)
	}
	if cuo.mutation.ErrorMessageCleared() {
		_spec.ClearField(check.FieldErrorMessage, field.TypeString)
	}
	if value, ok := cuo.mutation.HasDiff(); ok {
		_spec.SetField(check.FieldHasDiff, field.TypeBool, value)
	}
	if value, ok := cuo.mutation.DiffChange(); ok {
		_spec.SetField(check.FieldDiffChange, field.TypeJSON, value)
	}
	if cuo.mutation.DiffChangeCleared() {
		_spec.ClearField(check.FieldDiffChange, field.TypeJSON)
	}
	if value, ok := cuo.mutation.CreatedAt(); ok {
		_spec.SetField(check.FieldCreatedAt, field.TypeTime, value)
	}
	if cuo.mutation.WebsiteCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   check.WebsiteTable,
			Columns: []string{check.WebsiteColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(website.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.WebsiteIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   check.WebsiteTable,
			Columns: []string{check.WebsiteColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(website.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Check{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{check.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
