package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

var _ ent.Mixin = (*TenantID)(nil)

type TenantID struct{ mixin.Schema }

func (TenantID) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_id").
			Comment("租户id").
			NotEmpty().
			Unique().
			Immutable(),
	}
}

// Indexes of the AutoIncrementId.
func (TenantID) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
	}
}
