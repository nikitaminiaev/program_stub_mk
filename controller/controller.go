package controller

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (c *MkController) ProcessData() error {
	sensor, valueUint, ParamErr := c.getParam()

	if ParamErr != nil {
		return ParamErr
	}

	var err error
	switch sensor {
	case "servo_x":
		c.xCurrent = valueUint
		err = c.scanAlgorithmZ()
	case "servo_y":
		c.yCurrent = valueUint
		err = c.scanAlgorithmZ()
	case "servo_z":
		c.zCurrent = valueUint
		err = c.scanAlgorithmZ()
	default:
		err = errors.New("unknown sensor name")
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *MkController) getParam() (string, uint16, error) {
	data, open := <-c.inputCh

	if !open {
		return "", 0, errors.New("input channel is closed")
	}

	var dataMap map[string]any
	parsErr := json.Unmarshal([]byte(data), &dataMap)
	if parsErr != nil {
		return "", 0, parsErr
	}

	sensor, exists := dataMap["sensor"]
	if !exists {
		return "", 0, errors.New("sensor not exist")
	}
	sensorStr := ""
	if str, ok := sensor.(string); ok {
		sensorStr = str
	} else {
		return "", 0, errors.New("sensor in not string")
	}

	value, exists := dataMap["value"]
	if !exists {
		return "", 0, errors.New("value not exist")
	}

	var valueUint uint16
	if str, ok := value.(string); ok {
		temp, _ := strconv.ParseUint(str, 10, 16)
		valueUint = uint16(temp)
	} else if float, ok := value.(float64); ok {
		valueUint = uint16(float)
	} else {
		return "", 0, errors.New("value type is not wrong")
	}
	return sensorStr, valueUint, nil
}

func (c *MkController) scanAlgorithmZ() error {
	if c.yCurrent > surfaceSize || c.xCurrent > surfaceSize {
		return errors.New("coordinate > surfaceSize")
	}

	for z := c.zCurrent; z > 0; z-- {
		if z == c.SurfaceGenerator.Surface[c.yCurrent][c.xCurrent] {
			c.zCurrent = z + departureByZ
			c.outputCh <- fmt.Sprintf(`{"sensor": "surface", "z_val": %d}`, z)
			return nil
		}
	}

	c.zCurrent = departureByZ
	c.outputCh <- fmt.Sprintf(`{"sensor": "surface", "z_val": 0}`)
	return nil
}
