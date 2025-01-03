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
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/predicate"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/user"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/website"
	"github.com/gelleson/changescout/changescout/pkg/crons"
	"github.com/google/uuid"
)

// WebsiteUpdate is the builder for updating Website entities.
type WebsiteUpdate struct {
	config
	hooks    []Hook
	mutation *WebsiteMutation
}

// Where appends a list predicates to the WebsiteUpdate builder.
func (wu *WebsiteUpdate) Where(ps ...predicate.Website) *WebsiteUpdate {
	wu.mutation.Where(ps...)
	return wu
}

// SetName sets the "name" field.
func (wu *WebsiteUpdate) SetName(s string) *WebsiteUpdate {
	wu.mutation.SetName(s)
	return wu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableName(s *string) *WebsiteUpdate {
	if s != nil {
		wu.SetName(*s)
	}
	return wu
}

// SetURL sets the "url" field.
func (wu *WebsiteUpdate) SetURL(s string) *WebsiteUpdate {
	wu.mutation.SetURL(s)
	return wu
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableURL(s *string) *WebsiteUpdate {
	if s != nil {
		wu.SetURL(*s)
	}
	return wu
}

// SetCron sets the "cron" field.
func (wu *WebsiteUpdate) SetCron(ce crons.CronExpression) *WebsiteUpdate {
	wu.mutation.SetCron(ce)
	return wu
}

// SetNillableCron sets the "cron" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableCron(ce *crons.CronExpression) *WebsiteUpdate {
	if ce != nil {
		wu.SetCron(*ce)
	}
	return wu
}

// SetEnabled sets the "enabled" field.
func (wu *WebsiteUpdate) SetEnabled(b bool) *WebsiteUpdate {
	wu.mutation.SetEnabled(b)
	return wu
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableEnabled(b *bool) *WebsiteUpdate {
	if b != nil {
		wu.SetEnabled(*b)
	}
	return wu
}

// SetMode sets the "mode" field.
func (wu *WebsiteUpdate) SetMode(s string) *WebsiteUpdate {
	wu.mutation.SetMode(s)
	return wu
}

// SetNillableMode sets the "mode" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableMode(s *string) *WebsiteUpdate {
	if s != nil {
		wu.SetMode(*s)
	}
	return wu
}

// SetSetting sets the "setting" field.
func (wu *WebsiteUpdate) SetSetting(d *domain.Setting) *WebsiteUpdate {
	wu.mutation.SetSetting(d)
	return wu
}

// SetUserID sets the "user_id" field.
func (wu *WebsiteUpdate) SetUserID(u uuid.UUID) *WebsiteUpdate {
	wu.mutation.SetUserID(u)
	return wu
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableUserID(u *uuid.UUID) *WebsiteUpdate {
	if u != nil {
		wu.SetUserID(*u)
	}
	return wu
}

// ClearUserID clears the value of the "user_id" field.
func (wu *WebsiteUpdate) ClearUserID() *WebsiteUpdate {
	wu.mutation.ClearUserID()
	return wu
}

// SetNextCheckAt sets the "next_check_at" field.
func (wu *WebsiteUpdate) SetNextCheckAt(t time.Time) *WebsiteUpdate {
	wu.mutation.SetNextCheckAt(t)
	return wu
}

// SetNillableNextCheckAt sets the "next_check_at" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableNextCheckAt(t *time.Time) *WebsiteUpdate {
	if t != nil {
		wu.SetNextCheckAt(*t)
	}
	return wu
}

// SetLastCheckAt sets the "last_check_at" field.
func (wu *WebsiteUpdate) SetLastCheckAt(t time.Time) *WebsiteUpdate {
	wu.mutation.SetLastCheckAt(t)
	return wu
}

// SetNillableLastCheckAt sets the "last_check_at" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableLastCheckAt(t *time.Time) *WebsiteUpdate {
	if t != nil {
		wu.SetLastCheckAt(*t)
	}
	return wu
}

// ClearLastCheckAt clears the value of the "last_check_at" field.
func (wu *WebsiteUpdate) ClearLastCheckAt() *WebsiteUpdate {
	wu.mutation.ClearLastCheckAt()
	return wu
}

// SetCreatedAt sets the "created_at" field.
func (wu *WebsiteUpdate) SetCreatedAt(t time.Time) *WebsiteUpdate {
	wu.mutation.SetCreatedAt(t)
	return wu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (wu *WebsiteUpdate) SetNillableCreatedAt(t *time.Time) *WebsiteUpdate {
	if t != nil {
		wu.SetCreatedAt(*t)
	}
	return wu
}

// SetUpdatedAt sets the "updated_at" field.
func (wu *WebsiteUpdate) SetUpdatedAt(t time.Time) *WebsiteUpdate {
	wu.mutation.SetUpdatedAt(t)
	return wu
}

// SetUser sets the "user" edge to the User entity.
func (wu *WebsiteUpdate) SetUser(u *User) *WebsiteUpdate {
	return wu.SetUserID(u.ID)
}

// Mutation returns the WebsiteMutation object of the builder.
func (wu *WebsiteUpdate) Mutation() *WebsiteMutation {
	return wu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (wu *WebsiteUpdate) ClearUser() *WebsiteUpdate {
	wu.mutation.ClearUser()
	return wu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (wu *WebsiteUpdate) Save(ctx context.Context) (int, error) {
	wu.defaults()
	return withHooks(ctx, wu.sqlSave, wu.mutation, wu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (wu *WebsiteUpdate) SaveX(ctx context.Context) int {
	affected, err := wu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wu *WebsiteUpdate) Exec(ctx context.Context) error {
	_, err := wu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wu *WebsiteUpdate) ExecX(ctx context.Context) {
	if err := wu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wu *WebsiteUpdate) defaults() {
	if _, ok := wu.mutation.UpdatedAt(); !ok {
		v := website.UpdateDefaultUpdatedAt()
		wu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wu *WebsiteUpdate) check() error {
	if v, ok := wu.mutation.Name(); ok {
		if err := website.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Website.name": %w`, err)}
		}
	}
	if v, ok := wu.mutation.URL(); ok {
		if err := website.URLValidator(v); err != nil {
			return &ValidationError{Name: "url", err: fmt.Errorf(`ent: validator failed for field "Website.url": %w`, err)}
		}
	}
	if v, ok := wu.mutation.Cron(); ok {
		if err := website.CronValidator(string(v)); err != nil {
			return &ValidationError{Name: "cron", err: fmt.Errorf(`ent: validator failed for field "Website.cron": %w`, err)}
		}
	}
	return nil
}

func (wu *WebsiteUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := wu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(website.Table, website.Columns, sqlgraph.NewFieldSpec(website.FieldID, field.TypeUUID))
	if ps := wu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wu.mutation.Name(); ok {
		_spec.SetField(website.FieldName, field.TypeString, value)
	}
	if value, ok := wu.mutation.URL(); ok {
		_spec.SetField(website.FieldURL, field.TypeString, value)
	}
	if value, ok := wu.mutation.Cron(); ok {
		_spec.SetField(website.FieldCron, field.TypeString, value)
	}
	if value, ok := wu.mutation.Enabled(); ok {
		_spec.SetField(website.FieldEnabled, field.TypeBool, value)
	}
	if value, ok := wu.mutation.Mode(); ok {
		_spec.SetField(website.FieldMode, field.TypeString, value)
	}
	if value, ok := wu.mutation.Setting(); ok {
		_spec.SetField(website.FieldSetting, field.TypeJSON, value)
	}
	if value, ok := wu.mutation.NextCheckAt(); ok {
		_spec.SetField(website.FieldNextCheckAt, field.TypeTime, value)
	}
	if value, ok := wu.mutation.LastCheckAt(); ok {
		_spec.SetField(website.FieldLastCheckAt, field.TypeTime, value)
	}
	if wu.mutation.LastCheckAtCleared() {
		_spec.ClearField(website.FieldLastCheckAt, field.TypeTime)
	}
	if value, ok := wu.mutation.CreatedAt(); ok {
		_spec.SetField(website.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := wu.mutation.UpdatedAt(); ok {
		_spec.SetField(website.FieldUpdatedAt, field.TypeTime, value)
	}
	if wu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   website.UserTable,
			Columns: []string{website.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   website.UserTable,
			Columns: []string{website.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, wu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{website.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	wu.mutation.done = true
	return n, nil
}

// WebsiteUpdateOne is the builder for updating a single Website entity.
type WebsiteUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *WebsiteMutation
}

// SetName sets the "name" field.
func (wuo *WebsiteUpdateOne) SetName(s string) *WebsiteUpdateOne {
	wuo.mutation.SetName(s)
	return wuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableName(s *string) *WebsiteUpdateOne {
	if s != nil {
		wuo.SetName(*s)
	}
	return wuo
}

// SetURL sets the "url" field.
func (wuo *WebsiteUpdateOne) SetURL(s string) *WebsiteUpdateOne {
	wuo.mutation.SetURL(s)
	return wuo
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableURL(s *string) *WebsiteUpdateOne {
	if s != nil {
		wuo.SetURL(*s)
	}
	return wuo
}

// SetCron sets the "cron" field.
func (wuo *WebsiteUpdateOne) SetCron(ce crons.CronExpression) *WebsiteUpdateOne {
	wuo.mutation.SetCron(ce)
	return wuo
}

// SetNillableCron sets the "cron" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableCron(ce *crons.CronExpression) *WebsiteUpdateOne {
	if ce != nil {
		wuo.SetCron(*ce)
	}
	return wuo
}

// SetEnabled sets the "enabled" field.
func (wuo *WebsiteUpdateOne) SetEnabled(b bool) *WebsiteUpdateOne {
	wuo.mutation.SetEnabled(b)
	return wuo
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableEnabled(b *bool) *WebsiteUpdateOne {
	if b != nil {
		wuo.SetEnabled(*b)
	}
	return wuo
}

// SetMode sets the "mode" field.
func (wuo *WebsiteUpdateOne) SetMode(s string) *WebsiteUpdateOne {
	wuo.mutation.SetMode(s)
	return wuo
}

// SetNillableMode sets the "mode" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableMode(s *string) *WebsiteUpdateOne {
	if s != nil {
		wuo.SetMode(*s)
	}
	return wuo
}

// SetSetting sets the "setting" field.
func (wuo *WebsiteUpdateOne) SetSetting(d *domain.Setting) *WebsiteUpdateOne {
	wuo.mutation.SetSetting(d)
	return wuo
}

// SetUserID sets the "user_id" field.
func (wuo *WebsiteUpdateOne) SetUserID(u uuid.UUID) *WebsiteUpdateOne {
	wuo.mutation.SetUserID(u)
	return wuo
}

// SetNillableUserID sets the "user_id" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableUserID(u *uuid.UUID) *WebsiteUpdateOne {
	if u != nil {
		wuo.SetUserID(*u)
	}
	return wuo
}

// ClearUserID clears the value of the "user_id" field.
func (wuo *WebsiteUpdateOne) ClearUserID() *WebsiteUpdateOne {
	wuo.mutation.ClearUserID()
	return wuo
}

// SetNextCheckAt sets the "next_check_at" field.
func (wuo *WebsiteUpdateOne) SetNextCheckAt(t time.Time) *WebsiteUpdateOne {
	wuo.mutation.SetNextCheckAt(t)
	return wuo
}

// SetNillableNextCheckAt sets the "next_check_at" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableNextCheckAt(t *time.Time) *WebsiteUpdateOne {
	if t != nil {
		wuo.SetNextCheckAt(*t)
	}
	return wuo
}

// SetLastCheckAt sets the "last_check_at" field.
func (wuo *WebsiteUpdateOne) SetLastCheckAt(t time.Time) *WebsiteUpdateOne {
	wuo.mutation.SetLastCheckAt(t)
	return wuo
}

// SetNillableLastCheckAt sets the "last_check_at" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableLastCheckAt(t *time.Time) *WebsiteUpdateOne {
	if t != nil {
		wuo.SetLastCheckAt(*t)
	}
	return wuo
}

// ClearLastCheckAt clears the value of the "last_check_at" field.
func (wuo *WebsiteUpdateOne) ClearLastCheckAt() *WebsiteUpdateOne {
	wuo.mutation.ClearLastCheckAt()
	return wuo
}

// SetCreatedAt sets the "created_at" field.
func (wuo *WebsiteUpdateOne) SetCreatedAt(t time.Time) *WebsiteUpdateOne {
	wuo.mutation.SetCreatedAt(t)
	return wuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (wuo *WebsiteUpdateOne) SetNillableCreatedAt(t *time.Time) *WebsiteUpdateOne {
	if t != nil {
		wuo.SetCreatedAt(*t)
	}
	return wuo
}

// SetUpdatedAt sets the "updated_at" field.
func (wuo *WebsiteUpdateOne) SetUpdatedAt(t time.Time) *WebsiteUpdateOne {
	wuo.mutation.SetUpdatedAt(t)
	return wuo
}

// SetUser sets the "user" edge to the User entity.
func (wuo *WebsiteUpdateOne) SetUser(u *User) *WebsiteUpdateOne {
	return wuo.SetUserID(u.ID)
}

// Mutation returns the WebsiteMutation object of the builder.
func (wuo *WebsiteUpdateOne) Mutation() *WebsiteMutation {
	return wuo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (wuo *WebsiteUpdateOne) ClearUser() *WebsiteUpdateOne {
	wuo.mutation.ClearUser()
	return wuo
}

// Where appends a list predicates to the WebsiteUpdate builder.
func (wuo *WebsiteUpdateOne) Where(ps ...predicate.Website) *WebsiteUpdateOne {
	wuo.mutation.Where(ps...)
	return wuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (wuo *WebsiteUpdateOne) Select(field string, fields ...string) *WebsiteUpdateOne {
	wuo.fields = append([]string{field}, fields...)
	return wuo
}

// Save executes the query and returns the updated Website entity.
func (wuo *WebsiteUpdateOne) Save(ctx context.Context) (*Website, error) {
	wuo.defaults()
	return withHooks(ctx, wuo.sqlSave, wuo.mutation, wuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (wuo *WebsiteUpdateOne) SaveX(ctx context.Context) *Website {
	node, err := wuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (wuo *WebsiteUpdateOne) Exec(ctx context.Context) error {
	_, err := wuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wuo *WebsiteUpdateOne) ExecX(ctx context.Context) {
	if err := wuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wuo *WebsiteUpdateOne) defaults() {
	if _, ok := wuo.mutation.UpdatedAt(); !ok {
		v := website.UpdateDefaultUpdatedAt()
		wuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wuo *WebsiteUpdateOne) check() error {
	if v, ok := wuo.mutation.Name(); ok {
		if err := website.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Website.name": %w`, err)}
		}
	}
	if v, ok := wuo.mutation.URL(); ok {
		if err := website.URLValidator(v); err != nil {
			return &ValidationError{Name: "url", err: fmt.Errorf(`ent: validator failed for field "Website.url": %w`, err)}
		}
	}
	if v, ok := wuo.mutation.Cron(); ok {
		if err := website.CronValidator(string(v)); err != nil {
			return &ValidationError{Name: "cron", err: fmt.Errorf(`ent: validator failed for field "Website.cron": %w`, err)}
		}
	}
	return nil
}

func (wuo *WebsiteUpdateOne) sqlSave(ctx context.Context) (_node *Website, err error) {
	if err := wuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(website.Table, website.Columns, sqlgraph.NewFieldSpec(website.FieldID, field.TypeUUID))
	id, ok := wuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Website.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := wuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, website.FieldID)
		for _, f := range fields {
			if !website.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != website.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := wuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wuo.mutation.Name(); ok {
		_spec.SetField(website.FieldName, field.TypeString, value)
	}
	if value, ok := wuo.mutation.URL(); ok {
		_spec.SetField(website.FieldURL, field.TypeString, value)
	}
	if value, ok := wuo.mutation.Cron(); ok {
		_spec.SetField(website.FieldCron, field.TypeString, value)
	}
	if value, ok := wuo.mutation.Enabled(); ok {
		_spec.SetField(website.FieldEnabled, field.TypeBool, value)
	}
	if value, ok := wuo.mutation.Mode(); ok {
		_spec.SetField(website.FieldMode, field.TypeString, value)
	}
	if value, ok := wuo.mutation.Setting(); ok {
		_spec.SetField(website.FieldSetting, field.TypeJSON, value)
	}
	if value, ok := wuo.mutation.NextCheckAt(); ok {
		_spec.SetField(website.FieldNextCheckAt, field.TypeTime, value)
	}
	if value, ok := wuo.mutation.LastCheckAt(); ok {
		_spec.SetField(website.FieldLastCheckAt, field.TypeTime, value)
	}
	if wuo.mutation.LastCheckAtCleared() {
		_spec.ClearField(website.FieldLastCheckAt, field.TypeTime)
	}
	if value, ok := wuo.mutation.CreatedAt(); ok {
		_spec.SetField(website.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := wuo.mutation.UpdatedAt(); ok {
		_spec.SetField(website.FieldUpdatedAt, field.TypeTime, value)
	}
	if wuo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   website.UserTable,
			Columns: []string{website.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wuo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   website.UserTable,
			Columns: []string{website.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Website{config: wuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, wuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{website.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	wuo.mutation.done = true
	return _node, nil
}
