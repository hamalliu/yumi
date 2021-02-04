package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

// AliPayMerchant holds the schema definition for the AliPayMerchant entity.
type AliPayMerchant struct {
	ent.Schema
}

// Fields of the AliPayMerchant.
func (AliPayMerchant) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.String("seller_key").
			Unique(),
		field.String("app_id").
			NotEmpty(),
		field.String("public_key").
			NotEmpty(),
		field.String("private_key").
			NotEmpty(),
	}
}

// Edges of the AliPayMerchant.
func (AliPayMerchant) Edges() []ent.Edge {
	return nil
}
