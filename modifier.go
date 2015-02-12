package main

import "strings"

type replacer interface {
	Replace(s string) string
}

var escaper = strings.NewReplacer("/", "_", "?", "_", "\\", "_", ":", "_", "<", "_", ">", "_", "|", "_", "\"", "")
var expander = strings.NewReplacer(
	"Drms", "Drums",
	"Chrm", "Chromatic",
	"Snrs", "Snares",
	"Mrmba", "Marimba",
	"BBoo", "Bamboo",
	"Bambu", "Bamboo",
	"Whstle", "Whistle",
	"Whstl", "Whistle",
	"whistle", "Whistle",
	"Tweakr", "Tweaker",
	"Kloker", "Klocker",
	"Kalmba", "Kalimba",
	"Harms", "Harmonic",
	"Bng", "Bongo",
	"Trumpt", "Trumpet",
	"Skrach", "Skratch",
	"Ntrl", "Natural",
	"Ghatm", "Ghatham",
	"Segndo", "Segundo",
	"Seg", "Segundo",
	"Hmmr", "Hammer",
	"P.Domra", "Prima Domra",
	"Zhnghu", "Zhonghu",
	"Wurli", "Wurly",
	"Perc", "Percussion",
)
var markAndSpaceRemover = strings.NewReplacer("_", "", " ", "", "-", "", "+", "")

type nopReplacer struct{}

func (nr nopReplacer) Replace(s string) string {
	return s
}

type suffixRemover struct{}

func (sr suffixRemover) Replace(s string) string {
	s = escaper.Replace(s)
	s = strings.TrimSuffix(s, " w")
	s = strings.TrimSuffix(s, " mono")
	s = strings.TrimSuffix(s, " MW")
	return s
}

type wordExpander struct{}

func (we wordExpander) Replace(s string) string {
	s = suffixRemover{}.Replace(s)
	s = expander.Replace(s)
	s = markAndSpaceRemover.Replace(s)
	return s
}
