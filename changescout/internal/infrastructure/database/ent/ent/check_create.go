// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/check"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/website"
	"github.com/google/uuid"
)

// CheckCreate is the builder for creating a Check entity.
type CheckCreate struct {
	config
	mutation *CheckMutation
	hooks    []Hook
}

// SetWebsiteID sets the "website_id" field.
func (cc *CheckCreate) SetWebsiteID(u uuid.UUID) *CheckCreate {
	cc.mutation.SetWebsiteID(u)
	return cc
}

// SetNillableWebsiteID sets the "website_id" field if the given value is not nil.
func (cc *CheckCreate) SetNillableWebsiteID(u *uuid.UUID) *CheckCreate {
	if u != nil {
		cc.SetWebsiteID(*u)
	}
	return cc
}

// SetResult sets the "result" field.
func (cc *CheckCreate) SetResult(b []byte) *CheckCreate {
	cc.mutation.SetResult(b)
	return cc
}

// SetHasError sets the "has_error" field.
func (cc *CheckCreate) SetHasError(b bool) *CheckCreate {
	cc.mutation.SetHasError(b)
	return cc
}

// SetNillableHasError sets the "has_error" field if the given value is not nil.
func (cc *CheckCreate) SetNillableHasError(b *bool) *CheckCreate {
	if b != nil {
		cc.SetHasError(*b)
	}
	return cc
}

// SetErrorMessage sets the "error_message" field.
func (cc *CheckCreate) SetErrorMessage(s string) *CheckCreate {
	cc.mutation.SetErrorMessage(s)
	return cc
}

// SetNillableErrorMessage sets the "error_message" field if the given value is not nil.
func (cc *CheckCreate) SetNillableErrorMessage(s *string) *CheckCreate {
	if s != nil {
		cc.SetErrorMessage(*s)
	}
	return cc
}

// SetHasDiff sets the "has_diff" field.
func (cc *CheckCreate) SetHasDiff(b bool) *CheckCreate {
	cc.mutation.SetHasDiff(b)
	return cc
}

// SetNillableHasDiff sets the "has_diff" field if the given value is not nil.
func (cc *CheckCreate) SetNillableHasDiff(b *bool) *CheckCreate {
	if b != nil {
		cc.SetHasDiff(*b)
	}
	return cc
}

// SetDiffChange sets the "diff_change" field.
func (cc *CheckCreate) SetDiffChange(d *diff.Result) *CheckCreate {
	cc.mutation.SetDiffChange(d)
	return cc
}

// SetCreatedAt sets the "created_at" field.
func (cc *CheckCreate) SetCreatedAt(t time.Time) *CheckCreate {
	cc.mutation.SetCreatedAt(t)
	return cc
}

// SetID sets the "id" field.
func (cc *CheckCreate) SetID(u uuid.UUID) *CheckCreate {
	cc.mutation.SetID(u)
	return cc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (cc *CheckCreate) SetNillableID(u *uuid.UUID) *CheckCreate {
	if u != nil {
		cc.SetID(*u)
	}
	return cc
}

// SetWebsite sets the "website" edge to the Website entity.
func (cc *CheckCreate) SetWebsite(w *Website) *CheckCreate {
	return cc.SetWebsiteID(w.ID)
}

// Mutation returns the CheckMutation object of the builder.
func (cc *CheckCreate) Mutation() *CheckMutation {
	return cc.mutation
}

// Save creates the Check in the database.
func (cc *CheckCreate) Save(ctx context.Context) (*Check, error) {
	cc.defaults()
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CheckCreate) SaveX(ctx context.Context) *Check {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CheckCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CheckCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *CheckCreate) defaults() {
	if _, ok := cc.mutation.HasError(); !ok {
		v := check.DefaultHasError
		cc.mutation.SetHasError(v)
	}
	if _, ok := cc.mutation.HasDiff(); !ok {
		v := check.DefaultHasDiff
		cc.mutation.SetHasDiff(v)
	}
	if _, ok := cc.mutation.ID(); !ok {
		v := check.DefaultID()
		cc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *CheckCreate) check() error {
	if _, ok := cc.mutation.Result(); !ok {
		return &ValidationError{Name: "result", err: errors.New(`ent: missing required field "Check.result"`)}
	}
	if v, ok := cc.mutation.Result(); ok {
		if err := check.ResultValidator(v); err != nil {
			return &ValidationError{Name: "result", err: fmt.Errorf(`ent: validator failed for field "Check.result": %w`, err)}
		}
	}
	if _, ok := cc.mutation.HasError(); !ok {
		return &ValidationError{Name: "has_error", err: errors.New(`ent: missing required field "Check.has_error"`)}
	}
	if _, ok := cc.mutation.HasDiff(); !ok {
		return &ValidationError{Name: "has_diff", err: errors.New(`ent: missing required field "Check.has_diff"`)}
	}
	if _, ok := cc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Check.created_at"`)}
	}
	return nil
}

func (cc *CheckCreate) sqlSave(ctx context.Context) (*Check, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *CheckCreate) createSpec() (*Check, *sqlgraph.CreateSpec) {
	var (
		_node = &Check{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(check.Table, sqlgraph.NewFieldSpec(check.FieldID, field.TypeUUID))
	)
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cc.mutation.Result(); ok {
		_spec.SetField(check.FieldResult, field.TypeBytes, value)
		_node.Result = value
	}
	if value, ok := cc.mutation.HasError(); ok {
		_spec.SetField(check.FieldHasError, field.TypeBool, value)
		_node.HasError = value
	}
	if value, ok := cc.mutation.ErrorMessage(); ok {
		_spec.SetField(check.FieldErrorMessage, field.TypeString, value)
		_node.ErrorMessage = value
	}
	if value, ok := cc.mutation.HasDiff(); ok {
		_spec.SetField(check.FieldHasDiff, field.TypeBool, value)
		_node.HasDiff = value
	}
	if value, ok := cc.mutation.DiffChange(); ok {
		_spec.SetField(check.FieldDiffChange, field.TypeJSON, value)
		_node.DiffChange = value
	}
	if value, ok := cc.mutation.CreatedAt(); ok {
		_spec.SetField(check.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := cc.mutation.WebsiteIDs(); len(nodes) > 0 {
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
		_node.WebsiteID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CheckCreateBulk is the builder for creating many Check entities in bulk.
type CheckCreateBulk struct {
	config
	err      error
	builders []*CheckCreate
}

// Save creates the Check entities in the database.
func (ccb *CheckCreateBulk) Save(ctx context.Context) ([]*Check, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Check, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CheckMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *CheckCreateBulk) SaveX(ctx context.Context) []*Check {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CheckCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CheckCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
