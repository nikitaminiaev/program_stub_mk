package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

const (
	departureByZ = 10
	surfaceSize  = 76
)

type MkController struct {
	inputCh          chan string
	outputCh         chan string
	SurfaceGenerator SurfaceGenerator
	xCurrent         uint16
	yCurrent         uint16
	zCurrent         uint16
}

func NewController(inputCh chan string, outputCh chan string) MkController {
	return MkController{
		inputCh:          inputCh,
		outputCh:         outputCh,
		SurfaceGenerator: SurfaceGenerator{[surfaceSize][surfaceSize]uint16{}},
		xCurrent:         0,
		yCurrent:         0,
		zCurrent:         0,
	}
}

func (c *MkController) SetZ(z uint16) {
	c.zCurrent = z
}

func (c *MkController) GetX() uint16 {
	return c.xCurrent
}

func (c *MkController) GetY() uint16 {
	return c.yCurrent
}

func (c *MkController) GetZ() uint16 {
	return c.zCurrent
}

func (c *MkController) ProcessData() {
	data, open := <-c.inputCh

	if !open {
		return
	}

	var dataMap map[string]any
	err := json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		log.Fatal(err)

		return
	}

	sensor, exists := dataMap["sensor"]
	if !exists {
		log.Fatal("sensor not exist")

		return
	}

	value, exists := dataMap["value"]
	if !exists {
		log.Fatal("value not exist")
		return
	}

	var valueUint uint16
	if str, ok := value.(string); ok {
		temp, _ := strconv.ParseUint(str, 10, 16)
		valueUint = uint16(temp)
	} else if float, ok := value.(float64); ok {
		valueUint = uint16(float)
	} else {
		log.Fatal("value type is not wrong")
		return
	}

	switch sensor {
	case "servo_x":
		c.xCurrent = valueUint
		c.scanAlgorithmZ()
	case "servo_y":
		c.yCurrent = valueUint
		c.scanAlgorithmZ()
	case "servo_z":
		c.zCurrent = valueUint + 10
		c.scanAlgorithmZ()
	}
}

func (c *MkController) scanAlgorithmZ() {
	for z := c.zCurrent; z > 0; z-- {
		if z == c.SurfaceGenerator.Surface[c.yCurrent][c.xCurrent] {
			c.zCurrent = z + departureByZ
			c.outputCh <- fmt.Sprintf(`{"sensor": "surface", "z_val": %d}`, z)
			return
		}
	}

	c.zCurrent = departureByZ
	c.outputCh <- fmt.Sprintf(`{"sensor": "surface", "z_val": 0}`)
}
