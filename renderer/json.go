package renderer

import (
	"encoding/json"
	"github.com/fhoehl/sunnyday4pm/freezer"
)

func RenderJSON(icecream freezer.Icecream) string {
	doc, _ := json.Marshal(icecream)
	return string(doc)
}
