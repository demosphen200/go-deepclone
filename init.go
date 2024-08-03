package deepclone

import (
	"reflect"
	"time"
)

func init() {
	knownStructCloners = make(map[reflect.Type]CloneFn)
	RegisterStructCloner(
		func(src *time.Time) (*time.Time, error) {
			var tm = *src
			return &tm, nil
		},
	)
}
