package main

import (
	"bufio"
	"fmt"
	"os"
	muabdib "zentense/muabdib"
)

func main() {
	// Main loop listens to standard input
	fmt.Println("Muabdib v0.1")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if muabdib.ProcessUciCommand(line) {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		println("Exiting:", err.Error())
	}

}
