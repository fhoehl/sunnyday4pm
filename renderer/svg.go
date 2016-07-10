package renderer

import (
	"bytes"
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/fhoehl/sunnyday4pm/freezer"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math"
)

var (
	width  = float32(1080.0)
	height = float32(1080.0)
)

func mapValue64(value, min, max, toMin, toMax float64) float64 {
	return (value-min)*(toMax-toMin)/(max-min) + min
}

func mapValue(value, min, max, toMin, toMax float32) float32 {
	return (value-min)*(toMax-toMin)/(max-min) + toMin
}

func formatColor(color colorful.Color) string {
	return fmt.Sprintf(`style="fill:%s"`, color.Hex())
}

func formatRotate(angle float32, x float32, y float32) string {
	return fmt.Sprintf(`transform="rotate(%f, %f, %f)"`, angle, x, y)
}

func RenderSVG(icecream freezer.Icecream, buffer *bytes.Buffer) {
	canvas := svg.New(buffer)
	canvas.Startview(int(width), int(height), 0, 0, int(width), int(height))

	genes := icecream.Dna.Genes

	if len(genes) < 16 {
		log.Print("Empty genes", icecream.GenId, len(icecream.Dna.Genes))
		return
	}

	// Colors

	colors := icecream.Colors

	// Background

	h1, _, _ := colors[0].Hcl()
	h2, _, _ := colors[0].Hcl()
	h3, _, _ := colors[0].Hcl()
	backgroundH := math.Mod((h1+h2+h3)/3, 360.0)

	canvas.Rect(0, 0, int(width), int(height), formatColor(colorful.Hsv(backgroundH, .02, .98)))

	// Rect

	rectX := mapValue(genes[0], 0, 1, 0, width)
	rectY := mapValue(genes[1], 0, 1, 0, width)
	rectH := mapValue(genes[2], 0, 1, 0, width)
	rectW := mapValue(genes[3], 0, 1, 0, width)
	rectR := mapValue(genes[4], 0, 1, 0, 360)

	canvas.CenterRect(int(rectX), int(rectY), int(rectW), int(rectH), formatColor(colors[0]), formatRotate(rectR, rectX, rectY))

	// Circle

	circX := mapValue(genes[5], 0, 1, 0, width)
	circY := mapValue(genes[6], 0, 1, 0, width)
	circD := mapValue(genes[7], 0, 1, 0, width)

	canvas.Circle(int(circX), int(circY), int(circD/2), formatColor(colors[1]))

	// Triangle

	triX := mapValue(genes[8], 0, 1, 0, width)
	triY := mapValue(genes[9], 0, 1, 0, width)
	triH := mapValue(genes[10], 0, 1, 0, width)
	triW := mapValue(genes[11], 0, 1, 0, width)
	triR := mapValue(genes[12], 0, 1, 0, 360)

	triPointsX := make([]int, 4)
	triPointsY := make([]int, 4)

	triPointsX[0] = int(triX - triW/2)
	triPointsX[1] = int(triX)
	triPointsX[2] = int(triX + triW/2)
	triPointsX[3] = int(triX - triW/2)

	triPointsY[0] = int(triY - triH/2)
	triPointsY[1] = int(triY + triH/2)
	triPointsY[2] = int(triY - triH/2)
	triPointsY[3] = int(triY - triH/2)

	canvas.Polygon(triPointsX, triPointsY, formatColor(colors[2]), formatRotate(triR, triX, triY))

	canvas.End()
}
