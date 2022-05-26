package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/mixin"
	"github.com/google/uuid"
)

// Auth holds the schema definition for the Auth entity.
type Auth struct {
	ent.Schema
}

func (Auth) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (Auth) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("app_id", uuid.UUID{}),
		field.UUID("third_party_id", uuid.UUID{}),
		field.String("app_key"),
		field.String("app_secret"),
		field.String("redirect_url"),
	}
}

func (Auth) Edges() []ent.Edge {
	return nil
}

func (Auth) Indexes() []ent.Index {
	return []ent.Index{}
}
