package freezer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fhoehl/sunnyday4pm/db"
	"github.com/lucasb-eyer/go-colorful"
	"html/template"
	"log"
	"math/rand"
	"time"
)

var (
	client = db.GetInstance()
	r      = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Icecream struct {
	GenId     int
	Dna       Dna
	ColorDna  Dna
	Colors    [3]colorful.Color
	Parent1Id int
	Parent2Id int
}

type IcecreamSVG struct {
	Id       int
	Likes    float64
	Colors   [3]colorful.Color
	ColorDna Dna
	SvgDoc   template.HTML
}

type Dna struct {
	Genes []float32
}

type IcecreamError struct {
	s string
}

type ByGenId []IcecreamSVG

func (a ByGenId) Len() int           { return len(a) }
func (a ByGenId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGenId) Less(i, j int) bool { return a[i].Id > a[j].Id }

func (e *IcecreamError) Error() string {
	return e.s
}

func makeRandomColor(hMin, hMax, cMin, cMax, lMin, lMax float32) colorful.Color {
	return colorful.Hcl(
		float64(randomBetween(hMin, hMax)),
		float64(randomBetween(cMin, cMax)),
		float64(randomBetween(lMin, lMax)))
}

func (icecream *Icecream) makeColors() {
	genes := icecream.ColorDna.Genes
	idx := 0

	for i := range icecream.Colors {
		idx = i * 6
		hMin := mapValue(genes[idx], 0, 1, 0, 360)
		hMax := mapValue(genes[idx+1], 0, 1, 0, 360)
		cMin := mapValue(genes[idx+2], 0, 1, -1, 1)
		cMax := mapValue(genes[idx+3], 0, 1, -1, 1)
		lMin := genes[idx+4]
		lMax := genes[idx+5]

		icecream.Colors[i] = makeRandomColor(hMin, hMax, cMin, cMax, lMin, lMax)
	}
}

func MakeRandomIcecream() Icecream {
	genId := getNextGenerationId()
	dna := Dna{MakeRandomGenes(30)}
	colorDna := Dna{MakeRandomGenes(18)}
	var colors [3]colorful.Color

	icecream := Icecream{genId, dna, colorDna, colors, 0, 0}
	icecream.makeColors()

	return icecream
}

func LoadIcecreamById(id int) (icecream Icecream, err error) {
	key := Key(id)
	return LoadIcecream(key)
}

func LoadIcecream(key string) (icecream Icecream, err error) {
	icecreamJSON, err := client.Cmd("GET", key).Bytes()

	if err == nil {
		var icecream Icecream
		json.Unmarshal(icecreamJSON, &icecream)
		return icecream, nil
	}

	return MakeRandomIcecream(), &IcecreamError{"Could not find the ice cream"}
}

func (icecream Icecream) Save() {
	icecreamJson, _ := json.Marshal(icecream)

	key := icecream.Key()
	var err error

	// Save icecream data
	err = client.Cmd("SET", key, string(icecreamJson)).Err

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// Return a formatted Redis key for the give icecream
func (icecream Icecream) Key() string {
	return fmt.Sprintf(`icecream_%d`, icecream.GenId)
}

// Return a formatted Redis key for the given icecream id
func Key(id int) string {
	return fmt.Sprintf(`icecream_%d`, id)
}

// Apply a like to the given ice cream
func (icecream Icecream) Like() {
	icecream.Save()
}

func MakeRandomGenes(length int) []float32 {
	genes := make([]float32, length)

	for i := range genes {
		genes[i] = r.Float32()
	}

	return genes
}

func MakeRandomPopulation(size int, numberOfGenes int) {
	population := make([][]float32, size)
	for i := range population {
		population[i] = MakeRandomGenes(numberOfGenes)
	}
}

func Crossover(a Icecream, b Icecream, rate float32) (Icecream, error) {
	if len(a.Dna.Genes) != len(b.Dna.Genes) {
		return *new(Icecream), errors.New("Genes must be the same length")
	}

	crossover := MakeRandomIcecream()

	for i := range crossover.Dna.Genes {
		if r.Float32() < rate {
			crossover.Dna.Genes[i] = a.Dna.Genes[i]
		} else {
			crossover.Dna.Genes[i] = b.Dna.Genes[i]
		}
	}

	for x := 1; x < len(crossover.ColorDna.Genes); x += 2 {
		if r.Float32() < rate {
			crossover.ColorDna.Genes[x-1] = a.ColorDna.Genes[x-1]
			crossover.ColorDna.Genes[x] = a.ColorDna.Genes[x]
		} else {
			crossover.ColorDna.Genes[x-1] = b.ColorDna.Genes[x-1]
			crossover.ColorDna.Genes[x] = b.ColorDna.Genes[x]
		}
	}

	return crossover, nil
}

func (dna *Dna) Mutate(mutationRate float32) {
	for i := range dna.Genes {
		if r.Float32() < mutationRate {
			dna.Genes[i] = r.Float32()
		}
	}
}

func SelectParents() []Icecream {
	icecreams, err := client.Cmd("ZRANGEBYSCORE", "score", 1, "+inf").List()
	selectedIcecreams := make([]Icecream, 2)

	// Is there a least 2 parent?
	if err != nil || len(icecreams) < 2 {
		return make([]Icecream, 0)
	} else if len(icecreams) == 2 {
		selectedIcecreams[0], _ = LoadIcecream(icecreams[1])
		selectedIcecreams[1], _ = LoadIcecream(icecreams[0])
	} else {
		// Elite icecream is a the end of elements slice
		eliteIcecream := icecreams[len(icecreams)-1]

		// Randomly pick one
		pick := icecreams[r.Intn(len(icecreams)-1)]

		selectedIcecreams[0], _ = LoadIcecream(eliteIcecream)
		selectedIcecreams[1], _ = LoadIcecream(pick)
	}

	return selectedIcecreams
}

func MakeABaby() Icecream {
	var babyIcecream Icecream

	parents := SelectParents()

	if len(parents) == 0 {
		babyIcecream = MakeRandomIcecream()
		babyIcecream.Save()
	} else {
		log.Println("Elite parent", parents[0])
		log.Println("2nd parent", parents[1])
		babyIcecream, _ = Crossover(parents[0], parents[1], 0.80)
		babyIcecream.Parent1Id = parents[0].GenId
		babyIcecream.Parent2Id = parents[1].GenId
		log.Println("Baby", babyIcecream.ColorDna)
		babyIcecream.Dna.Mutate(.2)
		babyIcecream.ColorDna.Mutate(.2)
		log.Println("Mutated baby", babyIcecream.ColorDna)
		babyIcecream.makeColors()
		babyIcecream.Save()
		log.Println("Final baby", babyIcecream)
	}

	return babyIcecream
}

func getNextGenerationId() int {
	id, err := client.Cmd("INCR", "generation").Int()

	if err != nil {
		log.Fatal(err)
	}

	return id
}

func GetLastGenerationId() int {
	id, err := client.Cmd("GET", "generation").Int()

	if err != nil {
		log.Fatal(err)
	}

	return id
}
