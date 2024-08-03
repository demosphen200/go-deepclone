package deepclone

import (
	"errors"
	"reflect"
)

type CloneFn func(src reflect.Value, ptrToDst reflect.Value) error

var knownStructCloners map[reflect.Type]CloneFn

func RegisterStructCloner[T any](
	cloner func(src *T) (*T, error),
) {
	var ptrToT *T
	knownStructCloners[reflect.TypeOf(ptrToT).Elem()] = func(src reflect.Value, ptrToDst reflect.Value) error {
		value, ok := src.Interface().(T)
		if !ok {
			return errors.New("abnormal: src is not a valid type")
		}
		cloned, err := cloner(&value)
		rfCloned := reflect.ValueOf(cloned)
		if err != nil {
			return err
		}
		ptrToDst.Elem().Set(rfCloned.Elem())
		return nil
	}
}

func DeepCloneReflect(src reflect.Value, ptrToDst reflect.Value) error {
	//fmt.Printf("src=%s ptrToDst=%s\n", src.Type().String(), ptrToDst.Type().String())
	dst := ptrToDst
	if dst.Kind() != reflect.Ptr {
		return errors.New("ptrToDst is not a pointer")
	}
	if src.Type() != dst.Elem().Type() {
		return errors.New("src and dst have different types")
	}
	srcElem := src
	dstElem := dst.Elem()

	switch srcElem.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String, reflect.UnsafePointer:
		dstElem.Set(srcElem)
		return nil
	case reflect.Array:
		dstArray := reflect.New(srcElem.Type()).Elem()
		for index := 0; index < srcElem.Len(); index++ {
			err := DeepCloneReflect(srcElem.Index(index), dstArray.Index(index).Addr())
			if err != nil {
				return err
			}
		}
		dstElem.Set(dstArray)
	case reflect.Chan:
		dstChan := reflect.MakeChan(srcElem.Type(), srcElem.Cap())
		dstElem.Set(dstChan)
	case reflect.Func:
		if srcElem.Elem().Type() != dstElem.Elem().Type() {
			return errors.New("src and dst functions have different types")
		}
		dstElem.Set(srcElem)
	case reflect.Interface:
		if srcElem.IsNil() {
			dstElem.Set(srcElem)
		} else {
			ptrToSubDst := reflect.New(srcElem.Elem().Type())
			err := DeepCloneReflect(srcElem.Elem(), ptrToSubDst)
			if err != nil {
				return err
			}
			dstElem.Set(ptrToSubDst.Elem())
		}
	case reflect.Map:
		if srcElem.IsNil() {
			dstElem.Set(srcElem)
		} else {
			dstMap := reflect.MakeMapWithSize(srcElem.Type(), srcElem.Len())
			keys := srcElem.MapKeys()
			for _, key := range keys {
				currentSrcElem := srcElem.MapIndex(key)
				ptrToCurrentDstElem := reflect.New(currentSrcElem.Type())
				err := DeepCloneReflect(currentSrcElem, ptrToCurrentDstElem)
				if err != nil {
					return err
				}
				dstMap.SetMapIndex(key, ptrToCurrentDstElem.Elem())
			}
			dstElem.Set(dstMap)
		}
	case reflect.Pointer:
		if srcElem.IsNil() {
			dstElem.Set(srcElem)
		} else {
			ptrToSubDst := reflect.New(srcElem.Elem().Type())
			err := DeepCloneReflect(srcElem.Elem(), ptrToSubDst)
			if err != nil {
				return err
			}
			dstElem.Set(ptrToSubDst)
		}
	case reflect.Slice:
		if srcElem.IsNil() {
			dstElem.Set(srcElem)
		} else {
			dstSlice := reflect.MakeSlice(srcElem.Type(), srcElem.Len(), srcElem.Len())
			for index := 0; index < srcElem.Len(); index++ {
				err := DeepCloneReflect(srcElem.Index(index), dstSlice.Index(index).Addr())
				if err != nil {
					return err
				}
			}
			dstElem.Set(dstSlice)
		}
	case reflect.Struct:
		if cloner, found := knownStructCloners[srcElem.Type()]; found {
			err := cloner(srcElem, dst)
			if err != nil {
				return err
			}
		} else {
			for index := 0; index < srcElem.NumField(); index++ {
				srcField := srcElem.Field(index)
				if !srcField.CanSet() {
					continue
				}
				dstField := dstElem.Field(index)
				err := DeepCloneReflect(srcField, dstField.Addr())
				if err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("abnormal: unknown elem kind")
	}
	return nil
}

func DeepCloneAny(ptrToSrc any, ptrToDst any) error {
	src := reflect.ValueOf(ptrToSrc)
	dst := reflect.ValueOf(ptrToDst)
	if src.Kind() != reflect.Ptr {
		return errors.New("ptrToSrc is not a pointer")
	}
	if dst.Kind() != reflect.Ptr {
		return errors.New("ptrToDst is not a pointer")
	}
	if src.Elem().Type() != dst.Elem().Type() {
		return errors.New("src and dst have different types")
	}
	return DeepCloneReflect(reflect.ValueOf(ptrToSrc).Elem(), reflect.ValueOf(ptrToDst))
}

func DeepClone[T any](ptrToSrc *T, ptrToDst *T) error {
	return DeepCloneReflect(reflect.ValueOf(ptrToSrc).Elem(), reflect.ValueOf(ptrToDst))
}
