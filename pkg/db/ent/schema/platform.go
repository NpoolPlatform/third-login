package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/mixin"
	"github.com/google/uuid"
)

// Platform holds the schema definition for the Platform entity.
type Platform struct {
	ent.Schema
}

func (Platform) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (Platform) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("app_id", uuid.UUID{}),
		field.String("platform"),
		field.String("platform_auth_url"),
		field.String("logo_url"),
		field.String("platform_app_key"),
		field.String("platform_app_secret"),
		field.String("redirect_url"),
	}
}

func (Platform) Edges() []ent.Edge {
	return nil
}

func (Platform) Indexes() []ent.Index {
	return []ent.Index{}
}
