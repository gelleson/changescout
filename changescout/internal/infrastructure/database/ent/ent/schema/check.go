package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/google/uuid"
)

// Check holds the schema definition for the Check entity.
type Check struct {
	ent.Schema
}

// Fields of the Check.
func (Check) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("website_id", uuid.UUID{}).Optional(),
		field.Bytes("result").NotEmpty(),
		field.Bool("has_error").Default(false),
		field.String("error_message").Optional(),
		field.Bool("has_diff").Default(false),
		field.JSON("diff_change", &diff.Result{}).Optional(),
		field.Time("created_at"),
	}
}

// Edges of the Check.
func (Check) Edges() []ent.Edge {
	return []ent.Edge{
		// Website is the Website that this check belongs to.
		edge.To("website", Website.Type).
			Unique().
			Field("website_id"),
	}
}
