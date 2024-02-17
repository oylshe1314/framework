package http

import (
	"framework/errors"
	"mime/multipart"
	"reflect"
)

type FileHeaders map[string][]*multipart.FileHeader

func (files FileHeaders) read(v interface{}) error {
	var vt = reflect.TypeOf(v)
	if vt.Kind() != reflect.Pointer {
		return errors.Error("read get query: non-pointer")
	}

	vt = vt.Elem()
	if vt.Kind() != reflect.Struct {
		return errors.Error("read get query: non-struct")
	}

	var vv = reflect.ValueOf(v).Elem()
	var fn = vt.NumField()
	for i := 0; i < fn; i++ {
		var sf = vt.Field(i)
		var ft = sf.Type
		var name = sf.Tag.Get("file")
		if name == "-" {
			continue
		}

		if name == "" {
			name = sf.Name
		}

		var file = files[name]
		if len(file) == 0 {
			continue
		}

		var fl = len(file)
		var fv = vv.Field(i)
		switch ft.Kind() {
		case reflect.Pointer:
			ft = ft.Elem()
			if ft.Kind() != reflect.Struct {
				return errors.Error("")
			}

			if ft.PkgPath() != "mime/multipart" || ft.Name() != "FileHeader" {
				return errors.Error("")
			}

			fv.Set(reflect.ValueOf(file[0]))
		case reflect.Slice:
			var et = ft.Elem()
			if et.Kind() != reflect.Pointer {
				return errors.Error("")
			}

			et = et.Elem()
			if et.Kind() != reflect.Struct {
				return errors.Error("")
			}

			if et.PkgPath() != "mime/multipart" || et.Name() != "FileHeader" {
				return errors.Error("")
			}

			if fv.IsNil() {
				fv.Set(reflect.MakeSlice(ft, fl, fl))
			}

			for fi := 0; fi < fl; fi++ {
				fv.Index(fi).Set(reflect.ValueOf(file[fi]))
			}
		case reflect.Array:
			var et = ft.Elem()
			if et.Kind() != reflect.Pointer {
				return errors.Error("")
			}

			et = et.Elem()
			if et.Kind() != reflect.Struct {
				return errors.Error("")
			}

			if et.PkgPath() != "mime/multipart" || et.Name() != "FileHeader" {
				return errors.Error("")
			}

			if fl > fv.Len() {
				fl = fv.Len()
			}
			if fl == 0 {
				continue
			}

			for fi := 0; fi < fl; fi++ {
				fv.Index(fi).Set(reflect.ValueOf(file[fi]))
			}
		default:
			return errors.Error("")
		}
	}
	return nil
}
