package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"stubMk/controller"
	"testing"
)

func TestControllerProcessWithEmptySurfaceForIntOrStringData(t *testing.T) {
	type args struct {
		data string
	}

	tests := []struct {
		name   string
		args   args
		expect uint16
	}{
		{
			name:   "servo x = 1 int",
			args:   args{data: `{"sensor": "servo_x", "value": 1}`},
			expect: 1,
		},
		{
			name:   "servo x = 10 int",
			args:   args{data: `{"sensor": "servo_x", "value": 10}`},
			expect: 10,
		},
		{
			name:   "servo x = 55 int",
			args:   args{data: `{"sensor": "servo_x", "value": 55}`},
			expect: 55,
		},
		{
			name:   "servo x = 75 int",
			args:   args{data: `{"sensor": "servo_x", "value": 75}`},
			expect: 75,
		},
		{
			name:   "servo x = 1 string",
			args:   args{data: `{"sensor": "servo_x", "value": "1"}`},
			expect: 1,
		},
		{
			name:   "servo x = 10 string",
			args:   args{data: `{"sensor": "servo_x", "value": "10"}`},
			expect: 10,
		},
		{
			name:   "servo x = 55 string",
			args:   args{data: `{"sensor": "servo_x", "value": "55"}`},
			expect: 55,
		},
		{
			name:   "servo x = 75 string",
			args:   args{data: `{"sensor": "servo_x", "value": "75"}`},
			expect: 75,
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
			assert.Equal(t, tt.expect, c.GetX())
			var expected uint16 = 0
			assert.Equal(t, expected, c.GetY())

			close(inputCh)
			close(outputCh)
		})
	}
}

func TestControllerProcessWithFlatSurface(t *testing.T) {
	type args struct {
		data string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "servo x = 1",
			args: args{data: `{"sensor": "servo_x", "value": 1}`},
		},
		{
			name: "servo x = 10",
			args: args{data: `{"sensor": "servo_x", "value": 1}`},
		},
		{
			name: "servo x = 55",
			args: args{data: `{"sensor": "servo_x", "value": 55}`},
		},
		{
			name: "servo x = 75",
			args: args{data: `{"sensor": "servo_x", "value": 75}`},
		},
		{
			name: "servo y = 1",
			args: args{data: `{"sensor": "servo_y", "value": 1}`},
		},
		{
			name: "servo y = 10",
			args: args{data: `{"sensor": "servo_y", "value": 10}`},
		},
		{
			name: "servo y = 55",
			args: args{data: `{"sensor": "servo_y", "value": 55}`},
		},
		{
			name: "servo y = 75",
			args: args{data: `{"sensor": "servo_y", "value": 75}`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCh := make(chan string)
			outputCh := make(chan string)
			c := controller.NewController(inputCh, outputCh)
			var height uint16 = 11

			c.SurfaceGenerator.GenSurface(height)
			c.SetZ(12)

			go c.ProcessData()
			inputCh <- tt.args.data
			assert.Equal(t, fmt.Sprintf(`{"sensor": "surface", "z_val": %d}`, height), <-outputCh)
			assert.Equal(t, height, c.SurfaceGenerator.Surface[75][75])
			var expected uint16 = 21
			assert.Equal(t, expected, c.GetZ())

			close(inputCh)
			close(outputCh)
		})
	}

}
