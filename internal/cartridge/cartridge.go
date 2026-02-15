package cartridge

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// https://gbdev.io/pandocs/The_Cartridge_Header.html
type CartridgeHeader struct {
	// 0100-0103: entry point
	entry [4]byte
	// 0104-0133: nintendo logo
	logo [0x30]byte
	// 0134-0143 title in upper ASCII
	titleBytes [16]byte
	// 0144–0145: new licensee code
	newLicenseCode uint16
	// 0146: SGB flag
	sgbFlag byte
	// 0147: cartridge type
	cartType byte
	// 0148: ROM size
	romSize byte
	// 0149: RAM size
	ramSize byte
	// 014A: destination code (Japan or other)
	destCode byte
	// 014B: old licensee code
	oldLicenseeCode byte
	// 014C: ROM version number
	version byte
	// 014D: ROM checksum
	checksum byte
	// 014E-014F: global checksum
	// not used
	globalChecksum uint16
}

type CartridgeContext struct {
	romHeader CartridgeHeader
	romData []byte
	romSize uint32
	filename string
	title string
}

var ctx *CartridgeContext

func New(filename string) *CartridgeContext {
	return &CartridgeContext{
		filename: filename,
		romHeader: CartridgeHeader{},
	}
}

func buildCartridgeHeader(romData []byte) CartridgeHeader {
	return CartridgeHeader{
		titleBytes:      [16]byte(romData[0x134:0x144]),
		newLicenseCode:  uint16(romData[0x144])<<8 | uint16(romData[0x145]),
		cartType:        romData[0x147],
		romSize:         romData[0x148],
		ramSize:         romData[0x149],
		destCode:        romData[0x14A],
		oldLicenseeCode: romData[0x14B],
		version:         romData[0x14C],
		checksum:        romData[0x14D],
	}
}

func (c *CartridgeContext) getName() string {
	if c.title != "" {
		return c.title
	}

	c.romHeader.titleBytes = [16]byte{}
	for i := uint16(0x134); i <= 0x143; i++ {
		if c.romData[i] == 0x00 {
			break
		}
		c.romHeader.titleBytes[i-0x134] = c.romData[i]
	}

	c.title = strings.TrimSpace(string(c.romHeader.titleBytes[:]))
	return c.title
}

func Load(romPath string) error {
	ctx = New(romPath)

	file, err := os.Open(romPath)
	if err != nil {
		return fmt.Errorf("error opening rom file: %w", err)
	}
	defer file.Close()

	fmt.Println("rom file opened successfully")

	romData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading rom file: %w", err)
	}

	ctx.romData = romData
	ctx.romSize = uint32(len(romData))
	ctx.romHeader = buildCartridgeHeader(romData)

	fmt.Println()
	fmt.Println("--------------------------------")
	fmt.Println("cartridge loaded successfully!")
	fmt.Println("Title:", ctx.getName())
	fmt.Println("Type:", ROM_TYPES[ctx.romHeader.cartType])
	fmt.Println("ROM Size:", ROM_SIZES[ctx.romHeader.romSize])
	fmt.Println("RAM Size:", RAM_SIZES[ctx.romHeader.ramSize])
	fmt.Println("Destination Code:", DESTINATION_CODES[ctx.romHeader.destCode])
	newLicenseKey := string([]byte{byte(ctx.romHeader.newLicenseCode >> 8), byte(ctx.romHeader.newLicenseCode)})
	fmt.Println("Old Licensee Code:", OLD_LICENSEE_CODES[ctx.romHeader.oldLicenseeCode])
	fmt.Println("New Licensee Code:", NEW_LICENSEE_CODES[newLicenseKey])
	fmt.Println("Version:", ctx.romHeader.version)
	fmt.Println("--------------------------------")

	return nil
}
