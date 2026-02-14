package emulator

import (
	"fmt"
	"time"

	"github.com/matheus-foscarinid/gb-emu-go/internal/cartridge"
	"github.com/matheus-foscarinid/gb-emu-go/internal/cpu"
)

type EmulatorContext struct {
	paused bool
	running bool
	ticks uint64
}

var ctx *EmulatorContext

func New() *EmulatorContext {
	return &EmulatorContext{
		paused: false,
		running: false,
		ticks: 0,
	}
}

func Start(romPath string) error {
	ctx = New()
	ctx.running = true
	ctx.paused = false

	if err := cartridge.Load(romPath); err != nil {
		return err
	}

	fmt.Println("cartridge loaded successfully")

	return runLoop()
}

func runLoop() error {
	for ctx.running {
		if ctx.paused {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if err := cpu.Step(); err != nil {
			return fmt.Errorf("error on CPU step: %w", err)
		}

		ctx.ticks++
	}

	return nil
}
