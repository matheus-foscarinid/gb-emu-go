package cartridge

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	c := New("test.gb")
	if c == nil {
		t.Fatal("New returned nil")
	}
	if c.filename != "test.gb" {
		t.Errorf("filename = %q, want %q", c.filename, "test.gb")
	}
	if c.romData != nil {
		t.Error("romData should be nil for new context")
	}
	if c.romSize != 0 {
		t.Errorf("romSize = %d, want 0", c.romSize)
	}
}

func TestBuildCartridgeHeader(t *testing.T) {
	rom := minimalROM("HELLO WORLD")

	h := buildCartridgeHeader(rom)

	if got := string(h.titleBytes[:11]); got != "HELLO WORLD" {
		t.Errorf("titleBytes = %q, want HELLO WORLD", got)
	}
	if h.newLicenseCode[0] != '0' || h.newLicenseCode[1] != '1' {
		t.Errorf("newLicenseCode = %q, want \"01\"", string(h.newLicenseCode[:]))
	}
	if h.cartType != 0x00 {
		t.Errorf("cartType = 0x%02x, want 0x00", h.cartType)
	}
	if h.romSize != 0x00 {
		t.Errorf("romSize = 0x%02x, want 0x00", h.romSize)
	}
	if h.ramSize != 0x00 {
		t.Errorf("ramSize = 0x%02x, want 0x00", h.ramSize)
	}
	if h.destCode != 0x00 {
		t.Errorf("destCode = 0x%02x, want 0x00", h.destCode)
	}
	if h.oldLicenseeCode != 0x01 {
		t.Errorf("oldLicenseeCode = 0x%02x, want 0x01", h.oldLicenseeCode)
	}
	if h.version != 0x00 {
		t.Errorf("version = 0x%02x, want 0x00", h.version)
	}
	if h.checksum != rom[0x14D] {
		t.Errorf("checksum = 0x%02x, want 0x%02x", h.checksum, rom[0x14D])
	}
}

func TestGetName_FromROM(t *testing.T) {
	c := New("test.gb")
	c.romData = make([]byte, 0x144)
	copy(c.romData[0x134:0x144], "MY GAME")

	got := c.getName()
	if got != "MY GAME" {
		t.Errorf("getName() = %q, want %q", got, "MY GAME")
	}

	if c.title != "MY GAME" {
		t.Errorf("title not cached: %q", c.title)
	}
}

func TestGetName_Cached(t *testing.T) {
	c := New("test.gb")
	c.title = "Cached Title"

	got := c.getName()
	if got != "Cached Title" {
		t.Errorf("getName() = %q, want %q", got, "Cached Title")
	}
}

func TestGetName_StopsAtNull(t *testing.T) {
	c := New("test.gb")
	c.romData = make([]byte, 0x144)
	copy(c.romData[0x134:0x138], "AB\x00Z")

	got := c.getName()
	if got != "AB" {
		t.Errorf("getName() = %q, want %q", got, "AB")
	}
}

func minimalROM(title string) []byte {
	rom := make([]byte, 0x14E)
	if len(title) > 16 {
		title = title[:16]
	}
	copy(rom[0x134:0x134+len(title)], title)
	rom[0x144] = '0'
	rom[0x145] = '1'
	rom[0x147] = 0x00
	rom[0x148] = 0x00
	rom[0x149] = 0x00
	rom[0x14A] = 0x00
	rom[0x14B] = 0x01
	rom[0x14C] = 0x00

	var checksum int8
	for i := uint16(0x0134); i <= 0x014C; i++ {
		checksum = checksum - int8(rom[i]) - 1
	}
	rom[0x14D] = byte(checksum)
	return rom
}

func TestLoad_ValidMinimalROM(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "valid.gb")
	rom := minimalROM("LOAD TEST")
	if err := os.WriteFile(path, rom, 0644); err != nil {
		t.Fatalf("write test ROM: %v", err)
	}

	err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Global ctx was set by Load; basic sanity checks
	if ctx == nil {
		t.Fatal("ctx should be set after Load")
	}
	if ctx.romSize != uint32(len(rom)) {
		t.Errorf("ctx.romSize = %d, want %d", ctx.romSize, len(rom))
	}
	if ctx.romHeader.cartType != 0x00 {
		t.Errorf("ctx.romHeader.cartType = 0x%02x", ctx.romHeader.cartType)
	}
	if ctx.getName() != "LOAD TEST" {
		t.Errorf("getName() = %q, want LOAD TEST", ctx.getName())
	}
}

func TestLoad_ChecksumFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "badchecksum.gb")
	rom := minimalROM("BAD")
	rom[0x14D] ^= 0xFF // corrupt checksum

	if err := os.WriteFile(path, rom, 0644); err != nil {
		t.Fatalf("write test ROM: %v", err)
	}

	err := Load(path)
	if err == nil {
		t.Fatal("Load expected to fail on checksum error")
	}
	if err.Error() != "checksum validation failed" {
		t.Errorf("unexpected error: %v", err)
	}
}
