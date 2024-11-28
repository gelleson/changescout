package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/pkg/crons"
	"github.com/google/uuid"
	"time"
)

// Website holds the schema definition for the Website entity.
type Website struct {
	ent.Schema
}

// Fields of the Website.
func (Website) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			NotEmpty().
			MaxLen(255).
			Unique(),
		field.String("url").
			NotEmpty().
			MaxLen(255).
			Unique(),
		field.String("cron").
			NotEmpty().
			GoType(crons.CronExpression("")).
			MaxLen(255),
		field.Bool("enabled").
			Default(false),
		field.String("mode").
			Default("plain"),
		field.JSON("setting", &domain.Setting{}),
		field.UUID("user_id", uuid.UUID{}).
			Optional(),
		field.Time("next_check_at"),
		field.Time("last_check_at").
			Optional(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Website.
func (Website) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Field("user_id").Unique(),
	}
}

// Annotations of the Website.
func (Website) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "websites"},
	}
}

// Indexes of the Website.
func (Website) Indexes() []ent.Index {
	return []ent.Index{}
}
