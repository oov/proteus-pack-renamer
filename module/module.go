//go:generate go run gen.go -output modules.go

package module

// Module is a sound module.
type Module struct {
	Name  string
	Dir   string
	Banks []Bank
}

// Bank is a sound bank of the sound module.
type Bank struct {
	Name    string
	Dir     string
	Presets []Preset
}

// Preset is a sound preset of the sound module.
type Preset struct {
	Category string
	Name     string
}
