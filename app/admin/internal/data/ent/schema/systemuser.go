package schema

import (
	"qn-base/pkg/ent/mixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// SystemUser holds the schema definition for the SystemUser entity.
type SystemUser struct {
	ent.Schema
}

// Annotations of the User.
func (SystemUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_system_user"},
	}
}

// Fields of the SystemUser.
func (SystemUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("account").
			Unique().
			Comment("用户账号"),
		field.String("password").
			Optional().
			Nillable().
			Comment("密码"),
		field.String("nickname").
			Optional().
			Nillable().
			Comment("用户昵称"),
		field.String("remark").
			Optional().
			Nillable().
			Comment("备注"),
		field.String("dept_id").
			Optional().
			Nillable().
			Comment("部门ID"),
		field.String("post_ids").
			Optional().
			Nillable().
			Comment("岗位编号数组"),
		field.String("email").
			Optional().
			Nillable().
			Comment("用户邮箱"),
		field.String("mobile").
			Optional().
			Nillable().
			Comment("手机号码"),
		field.Int8("sex").
			Default(0).
			Optional().
			Nillable().
			Comment("用户性别(0:女 1:男)"),
		field.String("avatar").
			Optional().
			Nillable().
			Comment("头像地址"),
		field.Int8("status").
			Default(1).
			Optional().
			Comment("帐号状态(0:停用 1:正常)"),
		field.String("login_ip").
			Optional().
			Nillable().
			Comment("最后登录IP"),
		field.Time("login_date").
			Optional().
			Nillable().
			Comment("最后登录时间"),
	}
}

// Edges of the SystemUser.
func (SystemUser) Edges() []ent.Edge {
	return nil
}

// Mixin of the SystemUser.
func (SystemUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.StringId{},
		mixin.CreateBy{},
		mixin.CreateAt{},
		mixin.UpdateBy{},
		mixin.UpdateAt{},
		mixin.DeletedAt{},
		mixin.TenantID{},
	}
}
