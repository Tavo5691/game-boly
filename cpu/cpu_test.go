package cpu

import "testing"

func TestReset_8BitRegisters(t *testing.T) {
	testCases := []struct {
		desc     string
		register string
		testFn   func(*CPU) uint8
		expected uint8
	}{
		{
			desc:     "check A initial state",
			register: "A",
			testFn: func(c *CPU) uint8 {
				return c.GetA()
			},
			expected: 0x01,
		},
		{
			desc:     "check B initial state",
			register: "B",
			testFn: func(c *CPU) uint8 {
				return c.GetB()
			},
			expected: 0x00,
		},
		{
			desc:     "check C initial state",
			register: "C",
			testFn: func(c *CPU) uint8 {
				return c.GetC()
			},
			expected: 0x13,
		},
		{
			desc:     "check D initial state",
			register: "D",
			testFn: func(c *CPU) uint8 {
				return c.GetD()
			},
			expected: 0x00,
		},
		{
			desc:     "check E initial state",
			register: "E",
			testFn: func(c *CPU) uint8 {
				return c.GetE()
			},
			expected: 0xD8,
		},
		{
			desc:     "check F initial state",
			register: "F",
			testFn: func(c *CPU) uint8 {
				return c.GetF()
			},
			expected: 0xB0,
		},
		{
			desc:     "check H initial state",
			register: "H",
			testFn: func(c *CPU) uint8 {
				return c.GetH()
			},
			expected: 0x01,
		},
		{
			desc:     "check L initial state",
			register: "L",
			testFn: func(c *CPU) uint8 {
				return c.GetL()
			},
			expected: 0x4D,
		},
	}

	cpu := NewCPU()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			got := tC.testFn(cpu)

			if got != tC.expected {
				t.Errorf("%s register expected to have 0x%02x, but got 0x%02x", tC.register, tC.expected, got)
			}
		})
	}
}

func TestReset_16BitRegisters(t *testing.T) {
	testCases := []struct {
		desc     string
		register string
		testFn   func(*CPU) uint16
		expected uint16
	}{
		{
			desc:     "check SP initial state",
			register: "SP",
			testFn: func(c *CPU) uint16 {
				return c.GetSP()
			},
			expected: 0xFFFE,
		},
		{
			desc:     "check PC initial state",
			register: "PC",
			testFn: func(c *CPU) uint16 {
				return c.GetPC()
			},
			expected: 0x0100,
		},
	}

	cpu := NewCPU()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.testFn(cpu)

			if got != tC.expected {
				t.Errorf("%s register expected to have 0x%04x, but got 0x%04x", tC.register, tC.expected, got)
			}
		})
	}
}

func TestGettersAndSetters_ValuesIsTheExpectedOne(t *testing.T) {
	testCases := []struct {
		desc        string
		register    string
		testFn      func(*CPU) uint
		expected    uint
		is16BitPair bool
		expectedHi  uint8
		expectedLo  uint8
	}{
		{
			desc:     "",
			register: "A",
			testFn: func(c *CPU) uint {
				c.SetA(0xFF)

				return uint(c.GetA())
			},
			is16BitPair: false,
		},
	}

	cpu := NewCPU()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.testFn(cpu)

			if tC.expected != got {
				t.Errorf("Expected %s register to be 0x%04x but got: 0x%04x", tC.register, tC.expected, got)
			}

			if tC.is16BitPair {
				if tC.expectedHi != uint8(got)
			}
		})
	}
}
