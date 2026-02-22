package cpu

type AddrMode uint8
type InstructionType uint8

const (
	AddrImp AddrMode = iota // no operand
	AddrReg    // one register
	AddrRegReg // register <- register
	AddrRegD8  // register <- 8-bit immediate
	AddrRegD16 // register <- 16-bit immediate
	AddrRegMem // register <- memory[HL] (?)
	AddrMemReg // memory[HL] <- register
	AddrMemD8  // memory[HL] <- 8-bit immediate
	AddrMem    // operate on memory
	AddrD16    // jump to 16-bit address
	AddrD8     // jump relative by signed byte
)

const (
	InstNone InstructionType = iota
	InstNoop
	InstLoad
	InstIncrement
	InstDecrement
	InstAdd
	InstSubtract
)

type Instruction struct {
	Opcode uint8
	AddrMode AddrMode
	Operand uint16
}
