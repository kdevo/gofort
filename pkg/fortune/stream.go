package fortune

import (
	"bufio"
	"bytes"
)

type StreamFortuneTeller struct {
	*SimpleFortuneTeller
	// splitFunc is the function used for scanning the file.
	splitFunc bufio.SplitFunc
	// bufferSize is the buffer size to use.
	bufferSize uint32
	// maxBufferSize is the maximum buffer size that may be allocated while streaming.
	// If this is 0, the maxBufferSize limited to bufferSize.
	maxBufferSize uint32
}

func NewStreamFortuneTeller() *StreamFortuneTeller {
	return &StreamFortuneTeller{
		SimpleFortuneTeller: NewSimpleFortuneTeller(),
		maxBufferSize:       0,
		bufferSize:          8 * 1024,
		splitFunc: func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			idx := bytes.Index(data, []byte{'%', '\n'})
			if idx == 0 {
				return idx + 2, nil, nil
			} else if idx > 0 {
				return idx + 2, data[0:idx], nil
			}
			if atEOF {
				return len(data), data, nil
			}
			return 0, nil, nil
		},
	}
}

func (f *StreamFortuneTeller) WithMaxBufferSize(bs uint32) *StreamFortuneTeller {
	f.maxBufferSize = bs
	return f
}

func (f *StreamFortuneTeller) WithBufferSize(bs uint32) *StreamFortuneTeller {
	f.bufferSize = bs
	return f
}

func (f *StreamFortuneTeller) Fortune() string {
	// Reservoir sampling over all files
	reservoir := ""
	i := float32(1)

	buf := make([]byte, f.bufferSize)
	for _, path := range f.categories {
		file, err := texts.Open(path)
		if err != nil {
			f.err = err
			continue
		}
		scanner := bufio.NewScanner(file)
		scanner.Buffer(buf, int(f.maxBufferSize))
		scanner.Split(f.splitFunc)

		for ; scanner.Scan(); i++ {
			if f.rand.Float32() <= 1/i {
				reservoir = scanner.Text()
			}
		}
	}
	return reservoir
}
