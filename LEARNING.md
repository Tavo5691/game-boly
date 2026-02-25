# Learned Low-Level concepts from this project

## CPU

### Registers

We can think of registers as ultra-fast memory slots within the chip itself and are accessed by name (A, B, C, etc) rather than by memory address (0xC000).

The Game Boy CPU contains 7 general-purpose 8-bit registers, 1 "flag" special register and 2 16-bit registers.

The CPU can only do math on values it's already holding internally. So the typical workflow is:

1. Load a value from RAM into a register
2. Do math on it (add, subtract, compare, etc.)
3. Store the result back to RAM if needed

For example, to add two numbers stored in memory:
Load the first number into register A
Load the second into register B
Add B to A         ← result goes back into A
Store A back to RAM

The Game Boy's RAM is too slow to do math directly on — you always have to bring values into registers first.

#### Game Boy Registers

A ->  Accumulator — almost all arithmetic/logic results land here
B, C -> General purpose, often used as a pair (BC) for 16-bit values or loop counters
D, E -> General purpose, often used as a pair (DE) for memory addresses
H, L -> General purpose, but HL is special — it's the primary pointer register used to read/write memory
F -> Flags — not directly writable by programs; set automatically by arithmetic instructions
SP -> Stack Pointer — points to the top of the stack in memory
PC ->  Program Counter — address of the next instruction to execute

#### The F (Flags) Register

F holds 4 flags packed into its upper nibble. The lower 4 bits are always zero — the hardware ignores writes to them.

```
Bit 7: Z  (Zero)       — set when the result of an operation is 0
Bit 6: N  (Subtract)   — set when the last operation was a subtraction
Bit 5: H  (Half-carry) — set when there is a carry from bit 3 into bit 4
Bit 4: C  (Carry)      — set when there is a carry out of bit 7, or a borrow on subtraction
Bits 3-0: always 0
```

Example: if you write 0xFF to F, reading it back gives 0xF0 — the lower nibble is masked out.

#### 16-bit Register Pairs

Pairs of 8-bit registers can be read and written together as a single 16-bit value:

```
AF   (A is high byte, F is low byte)
BC   (B is high byte, C is low byte)
DE   (D is high byte, E is low byte)
HL   (H is high byte, L is low byte)
```

Example: if B=0x12 and C=0x34, then BC=0x1234.
Setting BC=0x1234 puts 0x12 into B and 0x34 into C.

The high/low byte order matters — getting it backwards is a common bug.

#### Initial State (post-boot)

The Game Boy has a boot ROM burned into the hardware that runs before the game cartridge. By the time the cartridge starts (at address 0x0100), the CPU is in this exact state:

```
A=0x01  F=0xB0  B=0x00  C=0x13
D=0x00  E=0xD8  H=0x01  L=0x4D
SP=0xFFFE  PC=0x0100
```

Emulators skip running the boot ROM and hardcode this state in a Reset() function instead.

---

## Current Progress

### Phase 0 — Project Skeleton ✓
- go.mod initialized (module: github.com/tavo5691/game-boly)
- All package stubs created and compiling
- Repo pushed to https://github.com/Tavo5691/game-boly

### Phase 1 — CPU Registers & Flags (in progress)
- Written: table-driven tests for 8-bit register initial state (cpu/cpu_test.go)
- Pending: tests for F mask, flag helpers (SetZ/GetZ etc.), 16-bit pairs, SP/PC initial state
- Pending: implementation in cpu/cpu.go
