package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/google/uuid"
	"time"
)

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

// Fields of the Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").NotEmpty(),
		field.String("type").GoType(domain.NotificationType(domain.TelegramNotificationType)),
		field.String("token").Optional(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("website_id", uuid.UUID{}).Optional().Nillable(),
		field.String("destination").Optional(),
		field.Time("created_at"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Notification.
func (Notification) Edges() []ent.Edge {
	return []ent.Edge{}
}
