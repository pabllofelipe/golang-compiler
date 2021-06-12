package main

import (
	"reflect"
	"testing"
)

func TestTOKEN_SCANNER(t1 *testing.T) {
	type fields struct {
		classe string
		lexema string
		tipo   string
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TOKEN
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TOKEN{
				classe: tt.fields.classe,
				lexema: tt.fields.lexema,
				tipo:   tt.fields.tipo,
			}
			if got := t.SCANNER(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("SCANNER() = %v, want %v", got, tt.want)
			}
		})
	}
}
