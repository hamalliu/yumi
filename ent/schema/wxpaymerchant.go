package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

// WxPayMerchant holds the schema definition for the WxPayMerchant entity.
type WxPayMerchant struct {
	ent.Schema
}

// Fields of the WxPayMerchant.
func (WxPayMerchant) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.String("seller_key").
			Unique(),
		field.String("app_id").
			NotEmpty(),
		field.String("mch_id").
			NotEmpty(),
		field.String("private_key").
			NotEmpty(),
		field.String("secret").
			Optional(),
	}
}

// Edges of the WxPayMerchant.
func (WxPayMerchant) Edges() []ent.Edge {
	return nil
}
