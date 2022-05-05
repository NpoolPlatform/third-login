package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/google/uuid"
)

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("content"),
		field.UUID("app_id", uuid.UUID{}),
		field.UUID("lang_id", uuid.UUID{}),
		field.String("used_for"),
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
		field.String("title"),
	}
}

func (Template) Edges() []ent.Edge {
	return nil
}

func (Template) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("app_id", "lang_id", "used_for").
			Unique(),
	}
}
