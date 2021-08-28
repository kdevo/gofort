package fortune

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type SimpleFortuneTeller struct {
	// categories maps category identifiers to file names.
	categories map[string]string
	// err stores the last occured error.
	err error
	// rand is a simple non-deterministic, non-cryptographic random generator.
	rand *rand.Rand
}

func NewSimpleFortuneTeller() *SimpleFortuneTeller {
	return &SimpleFortuneTeller{
		categories: categories(),
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ft *SimpleFortuneTeller) WithCategories(category ...string) *SimpleFortuneTeller {
	ft.categories = make(map[string]string)
	for _, c := range category {
		ft.categories[c] = fmt.Sprintf("%s/%s", embeddedDirectory, c)
	}
	return ft
}

func (f *SimpleFortuneTeller) Fortune() string {
	// Very naive implementation reading everything into memory
	var all []byte
	for _, path := range f.categories {
		file, err := texts.Open(path)
		if err != nil {
			f.err = err
			continue
		}

		content, err := io.ReadAll(file)
		if err != nil {
			f.err = err
			continue
		}
		all = append(all, content...)
	}
	split := bytes.Split(all, []byte{'%', '\n'})
	fortune := split[f.rand.Intn(len(split))]
	if len(fortune) == 0 {
		return f.Fortune()
	}
	return string(fortune)
}

func (f *SimpleFortuneTeller) Err() error {
	return f.err
}
