# Game Boy (DMG) Emulator in Go — TDD Plan

## Collaboration Workflow (per phase)

1. User writes test stubs for the current component
2. Claude reviews and points out missing edge cases or coverage gaps
3. User completes test assertions
4. User writes the implementation until tests pass
5. Claude reviews implementation for correctness and Go idioms
6. Move to next phase

---

## Project Directory Structure

```
game-boly/
├── go.mod
├── main.go                        # entry point (added in Phase 9)
├── cpu/
│   ├── cpu.go                     # registers, flags, Reset()
│   ├── cpu_test.go
│   ├── instructions.go            # opcode dispatch + implementations
│   ├── instructions_test.go
│   ├── cb_instructions.go         # CB-prefix opcodes
│   └── cb_instructions_test.go
├── memory/
│   ├── bus.go                     # Read/Write dispatch
│   ├── bus_test.go
│   └── regions.go                 # memory map constants
├── cartridge/
│   ├── cartridge.go               # header parsing, MBC interface
│   ├── cartridge_test.go
│   ├── mbc0.go                    # ROM-only
│   └── mbc1.go                    # ROM/RAM banking (Phase 10)
├── ppu/
│   ├── ppu.go                     # state machine, registers
│   ├── ppu_test.go
│   ├── renderer.go                # tile/sprite rendering
│   └── renderer_test.go
├── timer/
│   ├── timer.go
│   └── timer_test.go
├── interrupts/
│   ├── interrupts.go
│   └── interrupts_test.go
├── joypad/
│   ├── joypad.go
│   └── joypad_test.go
├── display/
│   └── display.go                 # Ebitengine wrapper (Phase 9)
└── testdata/
    ├── roms/                      # Blargg test ROMs (integration)
    └── fixtures/                  # golden frame buffers (PNG)
```

**Rule:** packages only import packages strictly below them in the stack. CPU knows nothing about PPU; bus knows nothing about CPU.

---

## MVP Phases

### Phase 0 — Project Skeleton ← START HERE
**Goal:** `go test ./...` passes; all packages compile.
- Initialize `go.mod`
- Create stub files and one placeholder test per package
- Establish test helpers (`newTestCPU`, `loadProgram`, etc.)

### Phase 1 — CPU Registers & Flags
**Goal:** Accurate register file and flag operations, fully tested.
- `cpu/cpu.go`: struct with A, B, C, D, E, H, L, F, SP, PC
- 16-bit pair getters/setters: AF, BC, DE, HL
- Flag helpers: SetZ/GetZ, SetN/GetN, SetH/GetH, SetC/GetC
- `Reset()` matching DMG post-boot state: A=0x01, F=0xB0, BC=0x0013, DE=0x00D8, HL=0x014D, SP=0xFFFE, PC=0x0100

### Phase 2 — Memory Bus + Cartridge Loader
**Goal:** Load a ROM; route reads/writes through the full memory map.
- `memory/bus.go`: `Read(uint16) uint8`, `Write(uint16, uint8)`
- Memory map: ROM 0x0000-0x7FFF, VRAM 0x8000-0x9FFF, WRAM 0xC000-0xDFFF, Echo 0xE000-0xFDFF, OAM 0xFE00-0xFE9F, prohibited 0xFEA0-0xFEFF (→0xFF), I/O 0xFF00-0xFF7F (stub), HRAM 0xFF80-0xFFFE, IE 0xFFFF
- `cartridge/cartridge.go`: header parsing (title, type, ROM size, checksum)
- `cartridge/mbc0.go`: ROM-only (type 0x00)

### Phase 3 — CPU Instruction Set (non-CB opcodes)
**Goal:** All 245 non-CB opcodes; executes a hand-crafted ROM loop.
- `cpu/instructions.go`: `[256]func(*CPU)` dispatch + `Step() int`
- CPU holds reference to Bus
- Order: NOP/HALT → 8-bit loads → 16-bit loads → 8-bit ALU → 16-bit ALU → jumps → calls/returns → rotates → misc
- `cycleTable [256]int` — assert cycle counts in every test

### Phase 4 — CB-Prefix Instructions
**Goal:** All 256 CB opcodes; Blargg cpu_instrs passes (integration).
- `cpu/cb_instructions.go`: RLC/RRC/RL/RR/SLA/SRA/SWAP/SRL, BIT, RES, SET
- Modify `Step()` to detect 0xCB and dispatch into CB table

### Phase 5 — Interrupts & Timer
**Goal:** Timer interrupt fires; HALT works; EI/DI/RETI correct.
- `interrupts/interrupts.go`: IE 0xFFFF, IF 0xFF0F; vectors VBlank 0x0040, STAT 0x0048, Timer 0x0050, Serial 0x0058, Joypad 0x0060
- `timer/timer.go`: DIV 0xFF04, TIMA 0xFF05, TMA 0xFF06, TAC 0xFF07; `Tick(cycles int)`
- Integrate interrupt check into `CPU.Step()` after each instruction

### Phase 6 — PPU Background Layer
**Goal:** Render a complete 160×144 background frame to a pixel buffer.
- `ppu/ppu.go`: LCDC 0xFF40, STAT 0xFF41, SCY/SCX, LY, LYC, BGP 0xFF47
- Mode state machine: OAM scan 80c → Drawing 172c → HBlank 204c; VBlank lines 144-153
- `ppu/renderer.go`: 2bpp tile decoding, tile map, scroll, palette
- PPU timing: 456 cycles/line, 70224 cycles/frame
- Golden-image integration test

### Phase 7 — PPU Sprites & Window Layer
**Goal:** Sprites and window complete the PPU.
- OAM: 40 sprites × 4 bytes; max 10/scanline
- Sprite attributes: flip X/Y, priority, OBP0/OBP1, 8×8 and 8×16
- Window: WY 0xFF4A, WX 0xFF4B; own internal line counter

### Phase 8 — Joypad
**Goal:** Readable button state; joypad interrupt.
- `joypad/joypad.go`: 0xFF00; direction nibble (bit 4) vs action nibble (bit 5); active-low

### Phase 9 — Display Integration
**Goal:** Real-time window at ~60fps.
- Library: **Ebitengine** (`github.com/hajimehoshi/ebiten/v2`) — pure Go, no CGO
- Main loop: ~70224 cycles/frame; CPU → PPU → Timer → Interrupts in lockstep

### Phase 10 — MBC1 (Larger ROMs)
**Goal:** ROM/RAM banking; unlocks most of the library.
- `cartridge/mbc1.go`: bank switching, RAM enable/disable, MBC1 bank-0 alias quirk

### Phase 11 (Optional) — APU
- 4 channels: Pulse 1 (sweep), Pulse 2, Wave, Noise. Start with Pulse 1.

---

## Blargg Test ROM Milestones

| ROM | Gates |
|-----|-------|
| `cpu_instrs` 03-09 | Phase 3 done |
| `instr_timing` | Phase 3 done |
| `cpu_instrs` 01, 10-11 | Phase 4 done |
| `cpu_instrs` 02 | Phase 5 done |
| `mem_timing` | Phase 5 done |

Integration tests: `//go:build integration` tag so `go test ./...` stays fast.

---

## Progress

- [x] Phase 0 — Project Skeleton
- [ ] Phase 1 — CPU Registers & Flags  ← IN PROGRESS
- [ ] Phase 2 — Memory Bus + Cartridge
- [ ] Phase 3 — CPU Instructions
- [ ] Phase 4 — CB Instructions
- [ ] Phase 5 — Interrupts & Timer
- [ ] Phase 6 — PPU Background
- [ ] Phase 7 — PPU Sprites & Window
- [ ] Phase 8 — Joypad
- [ ] Phase 9 — Display Integration
- [ ] Phase 10 — MBC1
- [ ] Phase 11 — APU (optional)
