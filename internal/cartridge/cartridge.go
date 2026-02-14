package cartridge

// https://gbdev.io/pandocs/The_Cartridge_Header.html
type Cartridge struct {
	// 0100-0103: entry point
	entry [4]byte
	// 0104-0133: nintendo logo
	logo [0x30]byte
	// 0134-0143 title in upper ASCII
	title [16]byte
	// 0144–0145: new licensee code
	newLicenseCode uint16
	// 0146: SGB flag
	sgbFlag uint8
	// 0147: cartridge type
	cartType uint8
	// 0148: ROM size
	romSize uint8
	// 0149: RAM size
	ramSize uint8
	// 014A: destination code (Japan or other)
	destCode uint8
	// 014B: old licensee code
	oldLicenseeCode uint8
	// 014C: ROM version number
	version uint8
	// 014D: ROM checksum
	checksum uint8
	// 014E-014F: global checksum
	// not used
	globalChecksum uint16
}

type CartridgeContext struct {
	romHeader Cartridge
	romData []byte
	romSize uint32
	filename string
}

func Load(romPath string) error {
	return nil
}
