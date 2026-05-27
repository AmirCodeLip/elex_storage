package shared_kernel

import (
	"fmt"
	"reflect"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/jinzhu/copier"
)

func MapToGrpc(dest interface{}, src interface{}) error {
	err := copier.CopyWithOption(dest, src, copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: reflect.TypeOf(""),          // string type
				DstType: reflect.TypeOf(uuid.UUID{}), // uuid.UUID type
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(string)
					if !ok {
						return nil, fmt.Errorf("expected string, got %T", src)
					}
					return uuid.Parse(s)
				},
			},
			{
				SrcType: reflect.TypeOf(""),                       // string type
				DstType: reflect.TypeOf((*uuid.UUID)(nil)).Elem(), // *uuid.UUID type
				Fn: func(src interface{}) (interface{}, error) {
					if src == nil {
						return nil, nil
					}

					s, ok := src.(string)
					if !ok || s == "" {
						return nil, nil
					}

					parsed, err := uuid.Parse(s)
					if err != nil {
						return nil, err
					}

					return &parsed, nil
				},
			},
			{
				SrcType: reflect.TypeOf(time.Time{}),
				DstType: reflect.TypeOf(int64(0)),
				Fn: func(src interface{}) (interface{}, error) {
					t, ok := src.(time.Time)
					if !ok {
						return nil, fmt.Errorf("expected time.Time, got %T", src)
					}
					return t.Unix(), nil
				},
			},
		},
	})
	return err
}

func MapFromGrpc(dest interface{}, src interface{}) error {
	err := copier.CopyWithOption(dest, src, copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: uuid.UUID{},
				DstType: "",
				Fn: func(src interface{}) (interface{}, error) {
					return src.(uuid.UUID).String(), nil
				},
			},
		},
	})
	return err
}

func TimeToUnix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}
