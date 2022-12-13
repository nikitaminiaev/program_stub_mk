package tests

import (
	"github.com/stretchr/testify/assert"
	"stubMk/controller"
	"testing"
)

func TestController_ProcessData_with_empty_surface(t *testing.T) {
	type args struct {
		data string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "servo_x_1",
			args: args{data: `{"sensor": "servo_x", "value": 1}`},
		},
		{
			name: "servo_x_10",
			args: args{data: `{"sensor": "servo_x", "value": 10}`},
		},
		{
			name: "servo_y_1",
			args: args{data: `{"sensor": "servo_y", "value": 1}`},
		},
		{
			name: "servo_y_10",
			args: args{data: `{"sensor": "servo_y", "value": 10}`},
		},
		{
			name: "servo_z_1",
			args: args{data: `{"sensor": "servo_z", "value": 1}`},
		},
		{
			name: "servo_z_10",
			args: args{data: `{"sensor": "servo_z", "value": 10}`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCh := make(chan string)
			outputCh := make(chan string)
			c := controller.NewController(inputCh, outputCh)

			go c.ProcessData()
			inputCh <- tt.args.data
			assert.Equal(t, `{"sensor": "surface", "z_val": 0}`, <-outputCh)

			close(inputCh)
			close(outputCh)
		})
	}
}

func TestController_ProcessData_with_flat_surface(t *testing.T) {
	type args struct {
		data string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "servo_x_12",
			args: args{data: `{"sensor": "servo_x", "value": 1}`},
		},
		{
			name: "servo_x_10",
			args: args{data: `{"sensor": "servo_x", "value": 10}`},
		},
		{
			name: "servo_y_1",
			args: args{data: `{"sensor": "servo_y", "value": 1}`},
		},
		{
			name: "servo_y_10",
			args: args{data: `{"sensor": "servo_y", "value": 10}`},
		},
		{
			name: "servo_z_12",
			args: args{data: `{"sensor": "servo_z", "value": 12}`},
		},
		{
			name: "servo_z_70",
			args: args{data: `{"sensor": "servo_z", "value": 70}`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCh := make(chan string)
			outputCh := make(chan string)
			c := controller.NewController(inputCh, outputCh)
			c.GenSurface(11)
			c.SetZ(12)
			go c.ProcessData()
			inputCh <- tt.args.data

			assert.Equal(t, `{"sensor": "surface", "z_val": 11}`, <-outputCh)

			close(inputCh)
			close(outputCh)
		})
	}
}
