package main

import (
	"fmt"
	"os"

	"github.com/matheus-foscarinid/gb-emu-go/internal/emulator"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("missing argument: gb-emu-go <rom_file>")
		os.Exit(1)
	}

	romPath := os.Args[1]
	err := emulator.Start(romPath)
	if err != nil {
		fmt.Println("error starting emulator:", err)
		os.Exit(1)
	}
}
