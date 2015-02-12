// +build ignore

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"go/format"

	"github.com/oov/proteus-pack-renamer/module"
)

var filename = flag.String("output", "modules.go", "output file name")

var program = `
package module

// Modules is sound modules.
var Modules = []Module{
	{{range .}}Module{
			Name: {{printf "%#v" .Name}},
			Dir: {{printf "%#v" .Dir}},
			Banks: []Bank{
				{{range $b := .Banks}}Bank{
						Name: {{printf "%#v" $b.Name}},
						Dir: {{printf "%#v" $b.Dir}},
						Presets: []Preset{
							{{range $p := .Presets}}Preset{Category: {{printf "%#v" .Category}}, Name: {{printf "%#v" .Name}}},
							{{end}}
						},
					},
				{{end}}
			},
		},
	{{end}}
}
`

func main() {
	flag.Parse()

	ml, err := buildModules()
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	t := template.Must(template.New("main").Parse(program))
	if err := t.Execute(&buf, ml); err != nil {
		log.Fatal(err)
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(*filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func buildModules() ([]*module.Module, error) {
	ml := []*module.Module{
		&module.Module{
			Name: "Proteus 2000",
			Dir:  "Proteus 2000",
			Banks: []module.Bank{
				module.Bank{
					Name: "Bank 0",
					Dir:  "Proteus 2000 Bank 0",
				},
				module.Bank{
					Name: "Bank 1",
					Dir:  "Proteus 2000 Bank 1",
				},
				module.Bank{
					Name: "Bank 2",
					Dir:  "Proteus 2000 Bank 2",
				},
				module.Bank{
					Name: "Bank 3",
					Dir:  "Proteus 2000 Bank 3",
				},
				module.Bank{
					Name: "Bank 4",
					Dir:  "Proteus 2000 Bank 4",
				},
				module.Bank{
					Name: "Bank 5",
					Dir:  "Proteus 2000 Bank 5",
				},
				module.Bank{
					Name: "Bank 6",
					Dir:  "Proteus 2000 Bank 6",
				},
				module.Bank{
					Name: "Bank 7",
					Dir:  "Proteus 2000 Bank 7",
				},
			},
		},
		&module.Module{
			Name: "Mo' Phatt",
			Dir:  "Mo Phatt",
			Banks: []module.Bank{
				module.Bank{
					Name: "Bank 0",
					Dir:  "Mo' Phatt Bank 0",
				},
				module.Bank{
					Name: "Bank 1",
					Dir:  "Mo' Phatt Bank 1",
				},
				module.Bank{
					Name: "Bank 2",
					Dir:  "Mo' Phatt Bank 2",
				},
				module.Bank{
					Name: "Bank 3",
					Dir:  "Mo' Phatt Bank 3",
				},
			},
		},
		&module.Module{
			Name: "Planet Earth",
			Dir:  "Planet Earth",
			Banks: []module.Bank{
				module.Bank{
					Name: "Bank 0",
					Dir:  "Planet Earth Bank 0",
				},
				module.Bank{
					Name: "Bank 1",
					Dir:  "Planet Earth Bank 1",
				},
				module.Bank{
					Name: "Bank 2",
					Dir:  "Planet Earth Bank 2",
				},
				module.Bank{
					Name: "Bank 3",
					Dir:  "Planet Earth Bank 3",
				},
			},
		},
		&module.Module{
			Name: "Vintage Pro",
			Dir:  "Vintage Pro",
			Banks: []module.Bank{
				module.Bank{
					Name: "Bank 0",
					Dir:  "Vintage Pro Bank 0",
				},
				module.Bank{
					Name: "Bank 1",
					Dir:  "Vintage Pro Bank 1",
				},
				module.Bank{
					Name: "Bank 2",
					Dir:  "Vintage Pro Bank 2",
				},
				module.Bank{
					Name: "Bank 3",
					Dir:  "Vintage Pro Bank 3",
				},
			},
		},
		&module.Module{
			Name: "Xtreme Lead",
			Dir:  "Xtreme Lead",
			Banks: []module.Bank{
				module.Bank{
					Name: "Bank 0",
					Dir:  "Xtreme Lead 1 Bank 0",
				},
				module.Bank{
					Name: "Bank 1",
					Dir:  "Xtreme Lead 1 Bank 1",
				},
				module.Bank{
					Name: "Bank 2",
					Dir:  "Xtreme Lead 1 Bank 2",
				},
				module.Bank{
					Name: "Bank 3",
					Dir:  "Xtreme Lead 1 Bank 3",
				},
			},
		},
	}

	for _, m := range ml {
		for i, b := range m.Banks {
			fn := fmt.Sprintf("%s/bank%d.csv", m.Dir, i)
			f, err := os.Open(fn)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			r := csv.NewReader(f)
			r.LazyQuotes = false
			for {
				line, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, fmt.Errorf("%s: %v", fn, err)
				}
				b.Presets = append(b.Presets, module.Preset{
					Category: line[1],
					Name:     line[2],
				})
			}
			m.Banks[i] = b
		}
	}
	return ml, nil
}
