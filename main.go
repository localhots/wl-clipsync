package main

import (
	"fmt"
	"flag"
	"time"
	"os/exec"
	"bytes"
	"log"
)


func main() {
	interval := flag.Duration("i", 100*time.Millisecond, "Clipboard polling interval")
	debug := flag.Bool("d", false, "Debug output")
	lotsofDebug := flag.Bool("dd", false, "Full debug")
	flag.Parse()

	fmt.Println("Starting clipboard sync")
	fmt.Println("Polling interval:", *interval)

	var unified string
	for range time.NewTicker(*interval).C {
		primary := getVal("-p")
		if *lotsofDebug {
			log.Println("Primary:", primary)
		}
		if primary != "" && primary != unified {
			unified = primary
			setVal(unified)
			if *debug {
				fmt.Println("Setting clip to", unified)
			}
			continue
		}

		selection := getVal()
		if *lotsofDebug {
			log.Println("Selection:", primary)
		}
		if selection != "" && selection != unified {
			unified = selection
			setVal(unified, "-p")
			if *debug {
				fmt.Println("Setting primary to", unified)
			}
		}
	}
}

func getVal(args ...string) string {
	args = append(args, "-n")
	out, _ := exec.Command("wl-paste", args...).Output()
	return string(out)
}

func setVal(val string, args... string) {
	cmd := exec.Command("wl-copy", args...)
	cmd.Stdin = bytes.NewBufferString(val)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to get paste buffer value (%v): %v", args, err)
		log.Printf("Value in question: %q", val)
	}
}
