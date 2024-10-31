package jsonselector

import (
	"errors"
	"strings"
)

// Errors
var (
	ErrInvalidSelectorFormat = errors.New("invalid selector format")
	ErrParentFieldNotAllowed = errors.New("parent field is not allowed")
	ErrFieldNotFound         = errors.New("field not found")
	ErrNotSupportedType      = errors.New("not supported type")
)

// Selector is an interface to find the field from the response body.
type Selector interface {
	Find(data map[string]interface{}, selector string) (string, error)
}

type basicSelector struct {
}

func (s *basicSelector) Find(data map[string]interface{}, selector string) (string, error) {
	if !strings.HasPrefix(selector, ".") {
		return "", ErrInvalidSelectorFormat
	}

	if strings.HasSuffix(selector, ".") {
		return "", ErrInvalidSelectorFormat
	}

	selectors := strings.Split(selector, ".")
	if len(selectors) < 2 {
		return "", ErrInvalidSelectorFormat
	}

	nested, ok := data[selectors[1]]
	if !ok {
		return "", ErrFieldNotFound
	}

	if len(selectors) == 2 {
		value, ok := nested.(string)
		if ok {
			return value, nil
		}

		return "", ErrNotSupportedType
	}

	next, ok := nested.(map[string]interface{})
	if !ok {
		return "", ErrParentFieldNotAllowed
	}

	return s.Find(next, "."+strings.Join(selectors[2:], "."))
}

// NewSelector creates a new selector instance to find the field from the response body.
func NewSelector() Selector {
	return &basicSelector{}
}
