package validator

import (
	"context"
	"errors"
	"github.com/go-inventory-v2/application/domain/external"
	"github.com/go-playground/validator/v10"
)

// Schema struct defines a validation schema for product requests.
type Schema struct {
    Validate func(context.Context, any) error
}

func (s *Schema) InventoryPatchSchema() Schema {
    return Schema{
        Validate: func(ctx context.Context, data any) error {
			
			req, ok := data.(external.ProductRequest)
			if !ok {
                return errors.New("schema validation failed ! Please check the request body and try again.")
            }

			if req.Inventory == nil {
                return errors.New("schema validation failed ! field inventory is mandatory")
            }

            return nil
        },
    }
}

// Use in ProductAdd and ProductPut.
func (s *Schema)ProductAddSchema() Schema {
    return Schema{
        Validate: func(ctx context.Context, data any) error {
            
			req, ok := data.(external.ProductRequest)
			if !ok {
                return errors.New("schema validation failed ! Please check the request body and try again.")
            }

			validate := validator.New()
			err := validate.Struct(req)
			if err != nil {
				return errors.New("schema validation failed ! Please check the request body data and try again.")
			}

            if req.Price == nil {
                return errors.New("schema validation failed ! field price is mandatory")
            }

            if req.Inventory == nil {
                return errors.New("schema validation failed ! field inventory is mandatory")
            }

            return nil
        },
    }
}
