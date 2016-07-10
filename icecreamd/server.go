package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fhoehl/sunnyday4pm/db"
	"github.com/fhoehl/sunnyday4pm/freezer"
	"github.com/fhoehl/sunnyday4pm/renderer"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	client = db.GetInstance()
	width  = float32(256.0)
	height = float32(256.0)
	r      = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func renderIcecream(iceCream freezer.Icecream) string {
	b := new(bytes.Buffer)
	renderer.RenderSVG(iceCream, b)
	s := strings.Split(b.String(), "<!-- Generated by SVGo -->")
	return s[1]
}

func editHandler(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("icecreamId")
	idString, _ := strconv.Atoi(id)
	icecream, _ := freezer.LoadIcecreamById(idString)
	icecream.Like()
	log.Println("Like", id)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	freezer.MakeABaby()

	icecreamKeys, _ := client.Cmd("KEYS", "icecream_*").List()
	icecreamsJSON, _ := client.Cmd("MGET", icecreamKeys).ListBytes()

	svgs := make([]freezer.IcecreamSVG, len(icecreamsJSON))

	for i, icecreamJSON := range icecreamsJSON {
		var _icecream freezer.Icecream

		json.Unmarshal(icecreamJSON, &_icecream)

		fitnessScore, _ := client.Cmd("ZRANK", "score", _icecream.Key()).Float64()

		svgs[i] = freezer.IcecreamSVG{
			_icecream.GenId,
			fitnessScore,
			_icecream.Colors,
			_icecream.ColorDna,
			template.HTML(renderIcecream(_icecream)),
		}
	}

	sort.Sort(freezer.ByGenId(svgs))

	t, _ := template.ParseFiles("./icecreamd/index.html")

	t.Execute(w, svgs)
}

func main() {
	http.Handle("/", http.HandlerFunc(indexHandler))
	http.Handle("/likes", http.HandlerFunc(editHandler))

	log.Println("Go to http://localhost:2003")

	err := http.ListenAndServe(":2003", nil)

	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
