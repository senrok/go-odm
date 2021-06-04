/*
@Author: Weny Xu
@Date: 2021/06/03 17:48
*/

package odm

import (
	"github.com/jinzhu/inflection"
	"reflect"
	"regexp"
	"strings"
)

// CollName returns a model's collection name. The `CollectionNameGetter` will be used
// if the model implements this interface. Otherwise, the collection name is inferred
// based on the model's type using reflection.
func CollName(m IModel) string {
	if collNameGetter, ok := m.(CollectionNameGetter); ok {
		return collNameGetter.CollectionName()
	}
	name := reflect.TypeOf(m).Elem().Name()

	return inflection.Plural(ToSnakeCase(name))
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase returns snake_case of the provided value.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
