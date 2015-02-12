package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/arbovm/levenshtein"

	"github.com/oov/proteus-pack-renamer/module"
)

type detectionType string

const (
	detectionExact          = detectionType(" ")
	detectionSemiExact      = detectionType("S")
	detectionSemiSemiExact  = detectionType("A")
	detectionManualApproved = detectionType("B")
	detectionManual         = detectionType("C")
	detectionFuzzy          = detectionType("D")
)

type instrumentFiles map[string]os.FileInfo

func (ifs *instrumentFiles) SortedKeys() []string {
	ss := make([]string, 0, len(*ifs))
	for k := range *ifs {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return ss
}

type similar struct {
	Score int
	Key   string
}

type similarSlice []similar

func (s similarSlice) Len() int           { return len(s) }
func (s similarSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s similarSlice) Less(i, j int) bool { return s[i].Score < s[j].Score }

func findMatchedPreset(p *module.Preset, insts instrumentFiles, src, dest replacer) string {
	n := src.Replace(p.Name)
	for k := range insts {
		if dest.Replace(k) == n {
			return k
		}
	}
	return ""
}

func findMostMatchedPreset(p *module.Preset, insts instrumentFiles, src, dest replacer) string {
	n := src.Replace(p.Name)
	ss := make(similarSlice, 0, len(insts))
	for k := range insts {
		ss = append(ss, similar{levenshtein.Distance(dest.Replace(k), n), k})
	}
	sort.Sort(ss)
	return ss[0].Key
}

func match(insts instrumentFiles, m module.Module, bankIndex int, presetIndex int, preset module.Preset) (detectionType, string) {
	if k, ok := manualDetection[fmt.Sprintf(`%s %d-%03d`, m.Name, bankIndex, presetIndex)]; ok {
		if k == "" {
			return detectionManualApproved, findMostMatchedPreset(&preset, insts, escaper, nopReplacer{})
		}
		return detectionManual, k
	}

	if k := findMatchedPreset(&preset, insts, escaper, nopReplacer{}); k != "" {
		return detectionExact, k
	}

	if k := findMatchedPreset(&preset, insts, suffixRemover{}, nopReplacer{}); k != "" {
		return detectionSemiExact, k
	}

	if k := findMatchedPreset(&preset, insts, wordExpander{}, markAndSpaceRemover); k != "" {
		return detectionSemiSemiExact, k
	}
	return detectionFuzzy, findMostMatchedPreset(&preset, insts, escaper, nopReplacer{})
}

func getInstMap(path, extension string) (instrumentFiles, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	insts := instrumentFiles{}
	for _, inst := range files {
		if inst.IsDir() || !strings.HasSuffix(inst.Name(), extension) {
			continue
		}
		insts[strings.TrimSuffix(inst.Name(), extension)] = inst
	}
	return insts, nil
}

func main() {
	basePath := flag.String("path", "", "path of Digital Sound Factory \"EMU_Proteus_Pack_Kontakt\" directory")
	executeRename := flag.Bool("rename", false, "execute rename")
	flag.Parse()
	if *basePath == "" {
		flag.Usage()
		return
	}

	for _, m := range module.Modules {
		for bankIndex, bank := range m.Banks {
			insts, err := getInstMap(filepath.Join(*basePath, m.Dir, bank.Dir), ".nki")
			if err != nil {
				log.Println(err, filepath.Join(*basePath, m.Dir, bank.Dir))
				continue
			}
			for pIndex, p := range bank.Presets {
				if len(insts) == 0 {
					fmt.Fprintf(os.Stderr, "X\t% -12s-%d\tXXX\t---:NOMORE\tthere are no more files\n", m.Name, bankIndex)
					break
				}

				dt, key := match(insts, m, bankIndex, pIndex, p)
				inst := insts[key]
				fmt.Printf("%s\t% -12s-%d\t%03d\t%s:%s\t%s\n", dt, m.Name, bankIndex, pIndex, p.Category, p.Name, inst.Name())
				if dt == detectionFuzzy {
					continue
				}
				if *executeRename {
					newName := fmt.Sprintf(`%03d-%s-%s`, pIndex, p.Category, inst.Name())
					err = os.Rename(
						filepath.Join(*basePath, m.Dir, bank.Dir, inst.Name()),
						filepath.Join(*basePath, m.Dir, bank.Dir, newName),
					)
					if err != nil {
						log.Println("could not rename:", err)
						continue
					}
				}
				delete(insts, key)
			}
			// generate unused instrument file list
			for _, k := range insts.SortedKeys() {
				fmt.Fprintf(os.Stderr, "x\t% -12s-%d\tXXX\t---:UNUSED\t%s\n", m.Name, bankIndex, insts[k].Name())
			}
		}
	}
}
