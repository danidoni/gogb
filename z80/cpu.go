package cpu

import (
  "fmt"
  "io/ioutil"
  "os"
)

var BIOS = [...]byte{
    0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
    0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
    0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
    0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
    0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
    0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
    0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
    0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
    0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xF2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
    0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
    0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
    0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
    0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
    0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3c, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x4C,
    0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x20, 0xFE, 0x23, 0x7D, 0xFE, 0x34, 0x20,
    0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x20, 0xFE, 0x3E, 0x01, 0xE0, 0x50,
  }

type MMU struct {
  Memory [65536]byte
}

func (mmu *MMU) Init() {
  for i := 0; i < 65536; i++ {
    mmu.Memory[i] = 0x00
  }
}

func (mmu *MMU) rw(address uint16) uint16 {
  return uint16(mmu.Memory[address+1]) << 8 + uint16(mmu.Memory[address])
}

func (mmu *MMU) pushIntoStack(address uint16, value uint16) {
  h := uint8(( value & 0xff00 ) / 255)
  l := uint8( value & 0x00ff )
  mmu.Memory[address-1] = l
  mmu.Memory[address-2] = h
}

func (mmu *MMU) LoadBios() {
  fmt.Printf("BIOS is %d bytes large\n", len(BIOS))

  for index,e := range BIOS {
    mmu.Memory[index] = e
  }
}

func (mmu *MMU) loadRom(fileName string) {
  byteArray, err := ioutil.ReadFile(fileName)
  if err != nil {
    fmt.Printf("Cannot open file name %s\n", fileName)
    os.Exit(1)
  }

  fmt.Printf("Rom is %d bytes large\n", len(byteArray))

  for i := 0; i < len(byteArray); i++ {
    mmu.Memory[i] = byteArray[i]
  }
}

type CPU struct {
  clock_m, clock_t int
  a, b, c, d, e, h, l uint8  // 8-bit registers
  Pc, sp uint16                    // 16-bit registers
  M, T int                      // Clock for last instruction
  f uint8
  Stop bool
  ops []interface{}

  MMU
}

func (cpu *CPU) Reset() {
  cpu.clock_m = 0
  cpu.clock_t = 0
  cpu.a = 0
  cpu.Pc = 0x0
  cpu.f = 0
  cpu.sp = 0x0
}

func (cpu *CPU) writeZeroFlag(val bool) {
  if val {
    cpu.f = cpu.f ^ ( 1 << 7 )
  }
}

func (cpu *CPU) writeSubstractionFlag(val bool) {
  if val {
    cpu.f = cpu.f ^ ( 1 << 6 )
  }
}

// TODO
/* When the sum of the least significant bits of the operands is avobe 9, it
 * cannot be represented in BCD, as it goes from A-D
 */
func (cpu *CPU) writeHalfcarryFlag(val bool) {
  if val {
    cpu.f = cpu.f ^ ( 1 << 5 )
  }
}

func (cpu *CPU) writeCarryFlag(val bool) {
  if val {
    cpu.f = cpu.f ^ ( 1 << 4 )
  }
}

func (cpu CPU) r_bc() uint16 {
  return (uint16(cpu.b) << 8 ) + uint16(cpu.c)
}

func (cpu *CPU) w_bc(bc uint16) {
  cpu.b = uint8( ( bc & 0xff00 ) / 0x100 )
  cpu.c = uint8( bc & 0xff )
}

func (cpu CPU) r_de() uint16 {
  return uint16( (cpu.d << 8 ) + cpu.e )
}

func (cpu *CPU) w_de(de uint16) {
  cpu.d = uint8( ( de & 0xff00 ) / 0x100 )
  cpu.e = uint8( de & 0xff )
}

func (cpu CPU) r_hl() uint16 {
  return uint16( (cpu.h << 8 ) + cpu.l )
}

func (cpu *CPU) w_hl(hl uint16) {
  cpu.h = uint8( ( hl & 0xff00 ) / 0x100 )
  cpu.l = uint8( hl & 0xff )
}

func (cpu *CPU) a16() uint16 {
  return uint16( ( cpu.a << 8 ) + cpu.f )
}

func (cpu *CPU) Dispatch(op byte) {
  if op == 0x00 {
    // NOP: Does nothing
    cpu.M = 1
    cpu.T = 4
    cpu.Pc++
  } else if op == 0x01 {
    // LD bc,d16 : Loads inmediate 2 byte data from into bc
    w := cpu.rw(cpu.Pc+1)
    cpu.w_bc(w)
    cpu.M = 3; cpu.T = 12
    cpu.Pc++
  } else if op == 0x02 {
    // LD (bc), a: Loads value from a into the memory position pointed by BC
    cpu.Memory[cpu.r_bc()] = cpu.a
    cpu.M = 1; cpu.T = 8
    cpu.Pc++
  } else if op == 0x03 {
    // INC bc: Increments bc and put result into bc
    cpu.w_bc(cpu.r_bc() + 1)
    cpu.M = 1; cpu.T = 8
    cpu.Pc++
  } else if op == 0x05 {
    // DEC b
    cpu.b = cpu.b - 1
    if cpu.b <= 1 {
      cpu.writeZeroFlag(true) 
    }
    cpu.writeSubstractionFlag(true)
    cpu.M = 1; cpu.T = 4
    cpu.Pc++
  } else if op == 0x13 {
    // INC de
    cpu.w_de(cpu.r_de() + 1)
    cpu.M = 1; cpu.T = 8
    cpu.Pc++
  } else if op == 0x1a {
    // LD a, (de)
    cpu.a = cpu.Memory[cpu.r_de()]
    cpu.M = 1; cpu.T = 8
    cpu.Pc++
  } else if op == 0x21 {
    // LD hl, d16: Load inmediate 2 byte data from d into sp
    w := cpu.rw(cpu.Pc+1)
    fmt.Printf("LD hl, %04x\n", w)
    cpu.w_hl(w)
    cpu.M = 3; cpu.T = 12
    cpu.Pc += 3
  } else if op == 0x22 {
    fmt.Printf("LD (hl+), a\n")
    cpu.Memory[cpu.r_hl()] = cpu.a
    cpu.w_hl(cpu.r_hl() + 1)
    cpu.M = 1; cpu.T = 8
    cpu.Pc++
  } else if op == 0x23 {
    fmt.Printf("INC hl\n")
    cpu.w_hl(cpu.r_hl() + 1)
    cpu.M = 1; cpu.T = 8
    cpu.Pc++
  } else if op == 0x31 {
    // LD sp, d16: Load inmediate 2 byte data from d into sp
    w := cpu.rw(cpu.Pc+1)
    fmt.Printf("LD sp, %04x\n", w)
    cpu.sp = w
    cpu.M = 3; cpu.T = 12
    cpu.Pc += 3
  } else if op == 0x32 {
    // LD (hl-), a: Load a into the address pointed by hl and then store the
    // decremented hl into hl
    fmt.Printf("LD (hl-), a\n")
    cpu.a = cpu.Memory[cpu.r_hl()]
    cpu.w_hl(cpu.r_hl() - 1)
    cpu.M = 1; cpu.T = 8
    cpu.Pc++
  } else if op == 0xaf {
    // XOR a: xor a and place it into a
    fmt.Printf("XOR %02x\n", cpu.a)
    cpu.a = cpu.a ^ cpu.a
    if cpu.a == 0 {
      cpu.f = cpu.f ^ ( 1 << 7 )
    }
    cpu.M = 1; cpu.T = 4
    cpu.Pc++
  } else if op == 0xc3 {
    address := uint16(cpu.Memory[cpu.Pc+1]) + (uint16(cpu.Memory[cpu.Pc+2]) << 8)
    fmt.Printf("JMP %04x\n", address)
    cpu.M = 3
    cpu.T = 16
    cpu.Pc = address
    // FIXME: Must push return address in stack
  } else if op == 0xcb { // Call a16 if Zero
    address := uint16(cpu.a)
    cpu.pushIntoStack(cpu.sp, cpu.Pc + 1)
    cpu.sp = cpu.sp - 2
    if cpu.isZero() {
      fmt.Printf("CALL %04x\n", address)
      cpu.Pc = address
    }
  } else if op == 0xcd { // Call a16
    address := cpu.a16()
    fmt.Printf("CALL %04x\n", address)
    cpu.M = 3; cpu.T = 16
    cpu.Pc = address
    cpu.sp -= 2
  } else if op == 0xfe {
    // CP d8: Compare inmediate byte data from d and a
    // FIXME: load INMEDIATE 8bit data via rb()
    fmt.Printf("CP a,%04x\n", cpu.d)
    if cpu.a == cpu.d {
      cpu.writeZeroFlag(true)
    }
    cpu.writeSubstractionFlag(true)
    if ( cpu.a - cpu.d ) < 0 {
      cpu.writeHalfcarryFlag(true)
    }
    if cpu.a < cpu.d {
      cpu.writeCarryFlag(true)
    }
    cpu.M = 2; cpu.T = 8
    cpu.Pc++
    // TODO: Must set flags!
  } else if op == 0xff {
    fmt.Printf("RST 0x38\n")
    cpu.Memory[cpu.sp] = uint8( ( cpu.Pc & 0xff00 ) / 0x100 )
    cpu.Memory[cpu.sp-1] = uint8( cpu.Pc & 0xff )
    cpu.Pc = 0x38
    cpu.M = 1; cpu.T = 16
    cpu.Pc++
  } else {
    cpu.Stop = true
  }
}

func (cpu CPU) DumpOp(op byte) {
  var mnemo string

  if op == 0x00 {
    mnemo = "NOP"
  } else if op == 0x05 {
    mnemo = "DEC b"
  } else if op == 0x13 {
    mnemo = "INC de"
  } else if op == 0x1a {
    mnemo = "LD a,(de)"
  } else {
    mnemo = "UNKNOWN OPCODE, stopping\n"
  }

  fmt.Printf("Pc:%08x op:%02x a:%02x b:%02x c:%02x d:%02x e:%02x h:%02x l:%02x sp:%04x f:%02x\t%s\n", cpu.Pc, op, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.h, cpu.l, cpu.sp, cpu.f, mnemo)
}

func (cpu *CPU) isZero() bool {
  return (cpu.f & 0x80) == 0x80
}

