package fortune

import (
	"embed"
	"fmt"
)

type FortuneTeller interface {
	// Fortune returns a random fortune.
	Fortune() string
	// Err returns an error in case any occurred when calling Fortune.
	Err() error
}

const embeddedDirectory = "texts"

//go:embed texts
var texts embed.FS

// categories returns a 1:1 mapping of available categories to full paths of the embedded text files.
func categories() map[string]string {
	categories := make(map[string]string)
	entries, _ := texts.ReadDir(embeddedDirectory)
	for _, entry := range entries {
		if !entry.IsDir() {
			categories[entry.Name()] = fmt.Sprintf("%s/%s", embeddedDirectory, entry.Name())
		}
	}
	return categories
}

// New returns the default FortuneTeller implementation.
func New() FortuneTeller {
	return NewStreamFortuneTeller()
}
