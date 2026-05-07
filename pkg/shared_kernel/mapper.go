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
				SrcType: "",
				DstType: uuid.UUID{},
				Fn: func(src interface{}) (interface{}, error) {
					return uuid.Parse(src.(string))
				},
			},
			{
				SrcType: "",
				DstType: (*uuid.UUID)(nil),
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
