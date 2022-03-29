package main

import (
	"fmt"
	"io"
	"os"

	"go.fergus.london/telemetry"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(
			"Usage: %s [filename]\n\n\t- [filename] file containing raw binary LTM binary data.\n",
			os.Args[0],
		)
		return
	}

	var (
		dataFile *os.File
		frames   []telemetry.DecodableFrame
		err      error
	)

	if dataFile, err = os.Open(os.Args[1]); err != nil {
		panic(err)
	}
	defer dataFile.Close()

	if frames, err = telemetry.Parse(dataFile); err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Printf("File '%s' contains %d frames.\n", os.Args[1], len(frames))
	if hasFlag("-v", "--verbose") {
		for _, f := range frames {
			fmt.Println(f)
		}
	}
}

func hasFlag(opts ...string) bool {
	for _, arg := range os.Args {
		for _, opt := range opts {
			if opt == arg {
				return true
			}
		}
	}

	return false
}
