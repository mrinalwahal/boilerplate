package todo

import (
	"context"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func Test_service_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx   context.Context
		title string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Todo
		wantErr bool
	}{
		{
			name: "Create Todo",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx:   context.Background(),
				title: "Test",
			},
			want:    &Todo{Title: "Test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				db: tt.fields.db,
			}
			got, err := s.Create(tt.args.ctx, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
