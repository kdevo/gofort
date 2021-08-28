package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// BenchmarkCompareNetBSDFortune compares with the "original" C version of fortune.
//
// Interpretation:
// Results highly depend on the distribution's packaged fortune version:
// - Versus the BSD fortune port on Alpine Linux, gofort is much slower.
//   This is likely due to the use of strfiles in the BSD port which make it easy to jump to a random fortune without reading a complete file.
// - Versus the fortune-mod distributed on Arch Linux, gofort is likely faster.
func BenchmarkCompareNetBSDFortune(b *testing.B) {
	// Resolve necessary benchmark binaries:
	fortunePath, err := fortunePath()
	if err != nil {
		b.Log(err)
		return
	}
	b.Logf("found `fortune`: %s", fortunePath)
	gofortPath, err := gofortPath()
	if err != nil {
		b.Log(err)
		return
	}
	b.Logf("found `%s`:  %s", NAME, gofortPath)

	benchmarks := []struct {
		name string
		path string
	}{
		{"gofort", gofortPath},
		{"fortune", fortunePath},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				command := exec.Command(bm.path)
				var out bytes.Buffer
				command.Stdout = &out
				err := command.Run()
				if err != nil {
					b.Fatalf("could not run: %v", err)
				}
			}
		})
	}

}

func gofortPath() (string, error) {
	globPaths := []string{fmt.Sprintf("../bin/%s-*-%s_%s", NAME, runtime.GOOS, runtime.GOARCH), fmt.Sprintf("../bin/%s*", NAME)}
	for _, globPath := range globPaths {
		matches, err := filepath.Glob(globPath)
		if len(matches) >= 1 && err == nil {
			return matches[0], nil
		}
	}
	return "", fmt.Errorf("%s not found. please build the binary before executing this benchmark.", NAME)
}

func fortunePath() (string, error) {
	fortunePath, err := exec.LookPath("fortune")
	if err != nil {
		return "", errors.New("fortune not found. please install before executing this benchmark.")
	}
	return fortunePath, nil
}
