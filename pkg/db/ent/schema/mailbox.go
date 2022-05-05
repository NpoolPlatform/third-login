package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/google/uuid"
)

// MailBox holds the schema definition for the MailBox entity.
type MailBox struct {
	ent.Schema
}

// Fields of the MailBox.
func (MailBox) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("app_id", uuid.UUID{}),
		field.UUID("from_user_id", uuid.UUID{}),
		field.UUID("to_user_id", uuid.UUID{}),
		field.Bool("already_read"),
		field.String("title"),
		field.String("content"),
		field.Uint32("create_at").
			DefaultFunc(func() uint32 {
				return uint32(time.Now().Unix())
			}),
		field.Uint32("update_at").
			DefaultFunc(func() uint32 {
				return uint32(time.Now().Unix())
			}).
			UpdateDefault(func() uint32 {
				return uint32(time.Now().Unix())
			}),
		field.Uint32("delete_at").
			DefaultFunc(func() uint32 {
				return 0
			}),
	}
}

// Edges of the MailBox.
func (MailBox) Edges() []ent.Edge {
	return nil
}
