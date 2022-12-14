package controller

type SurfaceGenerator struct {
	Surface [surfaceSize][surfaceSize]uint16
}

func (sg *SurfaceGenerator) GenSurface(height uint16) {
	for x, column := range sg.Surface {
		for y := range column {
			sg.Surface[x][y] = height
		}
	}
}
