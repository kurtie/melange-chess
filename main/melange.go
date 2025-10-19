package main

import (
	"bufio"
	"fmt"
	"os"
	melange "zentense/melange"
)

func main() {
	// Main loop listens to standard input
	fmt.Println("Melange v0.1")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if melange.ProcessUciCommand(line) {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		println("Exiting:", err.Error())
	}

}
