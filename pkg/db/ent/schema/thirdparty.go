package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/mixin"
	"github.com/google/uuid"
)

// ThirdParty holds the schema definition for the ThirdParty entity.
type ThirdParty struct {
	ent.Schema
}

func (ThirdParty) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (ThirdParty) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("brand_name"),
		field.String("logo"),
		field.String("domain"),
	}
}

func (ThirdParty) Edges() []ent.Edge {
	return nil
}

func (ThirdParty) Indexes() []ent.Index {
	return []ent.Index{}
}
