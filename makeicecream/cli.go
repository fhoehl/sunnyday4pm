package main

import (
	"bytes"
	"fmt"
	"github.com/fhoehl/sunnyday4pm/freezer"
	"github.com/fhoehl/sunnyday4pm/renderer"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	app            = kingpin.New("makeicecream", "A command-line for generating a new ice cream.")
	generate       = app.Command("generate", "Generate an ice cream")
	generateFormat = generate.Flag("format", "Output format (json|svg)").Default("svg").String()
	id             = app.Command("id", "Return last generated ice cream’s ID")
	render         = app.Command("render", "Render an ice cream")
	renderId       = render.Arg("id", "ID of the ice cream to render").Required().Int()
	renderFormat   = render.Flag("format", "Output format (json|svg)").Default("svg").String()
)

func main() {
	log.SetOutput(ioutil.Discard)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case generate.FullCommand():
		if *generateFormat == "json" || *generateFormat == "svg" {
			icecream := freezer.MakeABaby()

			if *generateFormat == "json" {
				fmt.Println(renderer.RenderJSON(icecream))
			} else if *generateFormat == "svg" {
				b := new(bytes.Buffer)
				renderer.RenderSVG(icecream, b)
				fmt.Println(b.String())
			}
		}

	case id.FullCommand():
		fmt.Println(freezer.GetLastGenerationId())

	case render.FullCommand():
		icecream, _ := freezer.LoadIcecreamById(int(*renderId))
		b := new(bytes.Buffer)
		renderer.RenderSVG(icecream, b)
		fmt.Println(b.String())
	}
}
