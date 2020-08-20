package main

import (
	"fmt"
	"flag"
	"time"
	"os/exec"
	"bytes"
)


func main() {
	interval := flag.Duration("i", 100*time.Millisecond, "Clipboard polling interval")
	flag.Parse()

	fmt.Println("Starting clipboard sync")
	fmt.Println("Polling interval:", *interval)

	var unified string
	for range time.NewTicker(*interval).C {
		primary := getVal("-p")
		if primary != unified {
			unified = primary
			setVal(unified)
			fmt.Println("Setting clip to", unified)
		}

		selection := getVal()
		if selection != unified {
			unified = selection
			setVal(unified, "-p")
			fmt.Println("Setting primary to", unified)
		}
	}
}

func getVal(args ...string) string {
	args = append(args, "-n")
	out, err := exec.Command("wl-paste", args...).Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

func setVal(val string, args... string) {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = bytes.NewBufferString(val)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}