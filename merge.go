package config

import "reflect"

func tryDefault(val, def any) any {
	if val == nil {
		return def
	}

	if reflect.ValueOf(val).IsZero() {
		return def
	}

	return val
}

// tryValueOrDefault tries to use default value if val is not defined.
func tryValueOrDefault[T any](val, def T) T {
	return tryDefault(val, def).(T)
}

// mergeLeft function merges the zero values of the left struct
// with corresponding values from the right struct.
// For types Chan, Func, Interface, Map, Pointer, Slice, UnsafePointer
// only the pointer is copied. Private fields will be ignored.
// Returns the merged struct.
func mergeLeft[T any](left, right T) T {
	lVal := reflect.ValueOf(left)
	rVal := reflect.ValueOf(right)

	if lVal.Kind() != reflect.Struct {
		return tryValueOrDefault(left, right)
	}

	mergedVal := reflect.New(reflect.TypeOf(left))
	if merged, ok := mergedVal.Interface().(*T); ok {
		*merged = left
	}
	mergedVal = reflect.Indirect(mergedVal)

	fields := lVal.NumField()
	for i := 0; i < fields; i++ {
		if !lVal.Type().Field(i).IsExported() {
			continue
		}

		lfield := lVal.Field(i)
		rfield := rVal.Field(i)
		mfield := mergedVal.Field(i)
		if !mfield.CanSet() {
			continue
		}

		var mergedField any
		if lfield.Kind() == reflect.Struct {
			mergedField = mergeLeft(lfield.Interface(), rfield.Interface())
		} else {
			mergedField = tryDefault(lfield.Interface(), rfield.Interface())
		}

		if mergedField == nil {
			continue
		}

		mergedVal := reflect.ValueOf(mergedField)
		mfield.Set(mergedVal)
	}

	return mergedVal.Interface().(T)
}
