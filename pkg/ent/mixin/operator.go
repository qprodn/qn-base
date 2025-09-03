package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

var _ ent.Mixin = (*CreateBy)(nil)

type CreateBy struct{ mixin.Schema }

func (CreateBy) Fields() []ent.Field {
	return []ent.Field{
		field.String("create_by").
			Comment("创建者ID").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*UpdateBy)(nil)

type UpdateBy struct{ mixin.Schema }

func (UpdateBy) Fields() []ent.Field {
	return []ent.Field{
		field.String("update_by").
			Comment("更新者ID").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*CreateAt)(nil)

type CreateAt struct{ mixin.Schema }

func (CreateAt) Fields() []ent.Field {
	return []ent.Field{
		// 创建时间
		field.Time("created_at").
			Comment("创建时间").
			Immutable().
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*UpdateAt)(nil)

type UpdateAt struct{ mixin.Schema }

func (UpdateAt) Fields() []ent.Field {
	return []ent.Field{
		// 更新时间
		field.Time("updated_at").
			Comment("更新时间").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*DeletedAt)(nil)

type DeletedAt struct{ mixin.Schema }

func (DeletedAt) Fields() []ent.Field {
	return []ent.Field{
		// 删除时间
		field.Time("deleted_at").
			Comment("删除时间").
			Optional().
			Nillable(),
	}
}
