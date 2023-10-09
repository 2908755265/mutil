package mcopy

import (
	"errors"
	"reflect"
)

var (
	ErrDestIsNotPointer = errors.New("the dest param is not a pointer")
	ErrNotSupport       = errors.New("only support string and number field")
)

func CopySameNamedField(dest any, src any) error {
	drt := reflect.TypeOf(dest)
	if drt.Kind() != reflect.Pointer {
		return ErrDestIsNotPointer
	}

	drv := reflect.ValueOf(dest).Elem()
	srt := reflect.TypeOf(src)
	srv := reflect.ValueOf(src)

	if srt.Kind() == reflect.Pointer {
		srt = srt.Elem()
		srv = srv.Elem()
	}
	for i := 0; i < srt.NumField(); i++ {
		dvf := drv.FieldByName(srt.Field(i).Name)
		if dvf.IsValid() {
			svf := srv.Field(i)
			if dvf.Kind() == svf.Kind() {
				err := setFieldVal(dvf, svf, dvf.Kind())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func setFieldVal(dest reflect.Value, src reflect.Value, t reflect.Kind) error {
	switch t {
	case reflect.String:
		dest.SetString(src.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dest.SetInt(src.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dest.SetUint(src.Uint())
	case reflect.Float64, reflect.Float32:
		dest.SetFloat(src.Float())
	default:
		return ErrNotSupport
	}
	return nil
}
