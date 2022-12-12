package tests

import (
	"github.com/stretchr/testify/assert"
	"stubMk/controller"
	"testing"
)

func TestController_ProcessData(t *testing.T) {
	type args struct {
		data string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{data: "{\"sensor\": \"servo_x\", \"value\": 1}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCh := make(chan string)
			outputCh := make(chan string)
			c := controller.NewController(inputCh, outputCh)
			//go func() {
			//	output := <-outputCh
			//	assert.Equal(t, "{{\"sensor\": \"surface\", \"z_val\": 0}}", output)
			//}()

			go c.ProcessData()
			inputCh <- tt.args.data
			assert.Equal(t, "{{\"sensor\": \"surface\", \"z_val\": 0}}", <-outputCh)

			close(inputCh)
			close(outputCh)
		})
	}
}
