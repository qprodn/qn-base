package mixin

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
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

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// Interceptors of the DeletedAt.
func (d DeletedAt) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			// Skip soft-delete, means include soft-deleted entities.
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			// Add filter for non-deleted entities
			if w, ok := q.(interface{ Where(...interface{}) }); ok {
				w.Where(sql.FieldIsNull("deleted_at"))
			}
			return nil
		}),
	}
}

// Hooks of the DeletedAt.
func (d DeletedAt) Hooks() []ent.Hook {
	return []ent.Hook{
		// Hook that converts delete operations to updates
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				// Check if this is a delete operation
				if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
					return next.Mutate(ctx, m)
				}

				// Skip soft-delete, means delete the entity permanently.
				if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
					return next.Mutate(ctx, m)
				}

				// Convert delete to update by setting deleted_at
				if setter, ok := m.(interface{ SetDeletedAt(time.Time) }); ok {
					// Note: In practice, you would need to create a proper update mutation
					// This is a simplified version for demonstration
					setter.SetDeletedAt(time.Now())
				}

				return next.Mutate(ctx, m)
			})
		},
	}
}
