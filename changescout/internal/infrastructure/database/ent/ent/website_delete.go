// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/predicate"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/website"
)

// WebsiteDelete is the builder for deleting a Website entity.
type WebsiteDelete struct {
	config
	hooks    []Hook
	mutation *WebsiteMutation
}

// Where appends a list predicates to the WebsiteDelete builder.
func (wd *WebsiteDelete) Where(ps ...predicate.Website) *WebsiteDelete {
	wd.mutation.Where(ps...)
	return wd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wd *WebsiteDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, wd.sqlExec, wd.mutation, wd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (wd *WebsiteDelete) ExecX(ctx context.Context) int {
	n, err := wd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wd *WebsiteDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(website.Table, sqlgraph.NewFieldSpec(website.FieldID, field.TypeUUID))
	if ps := wd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, wd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	wd.mutation.done = true
	return affected, err
}

// WebsiteDeleteOne is the builder for deleting a single Website entity.
type WebsiteDeleteOne struct {
	wd *WebsiteDelete
}

// Where appends a list predicates to the WebsiteDelete builder.
func (wdo *WebsiteDeleteOne) Where(ps ...predicate.Website) *WebsiteDeleteOne {
	wdo.wd.mutation.Where(ps...)
	return wdo
}

// Exec executes the deletion query.
func (wdo *WebsiteDeleteOne) Exec(ctx context.Context) error {
	n, err := wdo.wd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{website.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wdo *WebsiteDeleteOne) ExecX(ctx context.Context) {
	if err := wdo.Exec(ctx); err != nil {
		panic(err)
	}
}
