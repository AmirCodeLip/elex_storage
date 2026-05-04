package shared_kernel

import (
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
