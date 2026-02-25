package cpu

type CPU struct {
	a, b, c, d, e, h, l, f uint8
	sp, pc                 uint16
}

func NewCPU() *CPU {
	return Reset()
}

func Reset() *CPU {
	return &CPU{}
}

// 8-bit Getters
func (c *CPU) GetA() uint8 {
	return c.a
}

func (c *CPU) GetB() uint8 {
	return c.b
}

func (c *CPU) GetC() uint8 {
	return c.c
}

func (c *CPU) GetD() uint8 {
	return c.d
}

func (c *CPU) GetE() uint8 {
	return c.e
}

func (c *CPU) GetF() uint8 {
	return c.f
}

func (c *CPU) GetH() uint8 {
	return c.h
}

func (c *CPU) GetL() uint8 {
	return c.l
}

// 16-bit Getters
func (c *CPU) GetSP() uint16 {
	return c.sp
}

func (c *CPU) GetPC() uint16 {
	return c.pc
}

// Setters
func (c *CPU) SetA(val uint8) {
	c.a = val
}
