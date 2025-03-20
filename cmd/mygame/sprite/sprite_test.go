package sprite_test

import (
	"reflect"
	"testing"

	"github.com/dwethmar/beherit/cmd/mygame/sprite"
)

func TestNewSpritesheet(t *testing.T) {
	t.Run("all sprites should be created", func(t *testing.T) {
		// read all properties from the sprite sheet
		s, err := sprite.NewSpritesheet()
		if err != nil {
			t.Error(err)
		}

		// loop through all properties of the sprite sheet with reflection
		// and check if they are not nil
		// if they are nil, the sprite sheet is not properly initialized
		v := reflect.ValueOf(s).Elem() // Get the underlying value of the pointer to Sprites
		r := v.Type()                  // Get the type of the struct

		for i := range v.NumField() {
			fieldValue := v.Field(i)
			fieldName := r.Field(i).Name

			if fieldValue.IsNil() {
				t.Errorf("Field %s is nil", fieldName)
			}
		}
	})
}
