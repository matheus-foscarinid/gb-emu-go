package cpu

import (
	"fmt"

	"github.com/matheus-foscarinid/gb-emu-go/internal/bus"
)

type CPUContext struct {
	A, F, B, C, D, E, H, L uint8
	SP, PC                  uint16
	Halted                  bool
	CurrentInstruction      Instruction
}

func (c *CPUContext) fetchInstruction() error {
	opcode, err := bus.Read(c.PC)
	if err != nil {
		return fmt.Errorf("error reading opcode: %w", err)
	}

	c.CurrentInstruction.Opcode = opcode
	return nil
}

func (c *CPUContext) fetchData() {}

func (c *CPUContext) executeInstruction() {}

func (c *CPUContext) Step() {
	if c.Halted {
			return
	}

	c.fetchInstruction()
	c.fetchData()
	c.executeInstruction()
}
