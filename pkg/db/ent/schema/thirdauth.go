package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/mixin"
	"github.com/google/uuid"
)

// ThirdAuth holds the schema definition for the ThirdAuth entity.
type ThirdAuth struct {
	ent.Schema
}

func (ThirdAuth) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (ThirdAuth) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("app_id", uuid.UUID{}),
		field.String("third"),
		field.String("logo_url"),
		field.String("third_app_key"),
		field.String("third_app_secret"),
		field.String("redirect_url"),
	}
}

func (ThirdAuth) Edges() []ent.Edge {
	return nil
}

func (ThirdAuth) Indexes() []ent.Index {
	return []ent.Index{}
}
