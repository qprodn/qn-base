package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type StringId struct {
	mixin.Schema
}

func (StringId) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Comment("id").
			NotEmpty().
			Unique().
			Immutable(),
	}
}

// Indexes of the StringId.
func (StringId) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id"),
	}
}
