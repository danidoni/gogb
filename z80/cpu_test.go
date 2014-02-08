package cpu

import (
  "testing"
  "fmt"
)

func TestResetCPU(t *testing.T) {
  cpu := CPU{}
  cpu.Reset()

  if cpu.clock_m != 0 {
    t.Error("m clock should be 0")
  }
  if cpu.clock_t != 0 {
    t.Error("t clock should be 0")
  }
  if cpu.a != 0 {
    t.Error("a register should be 0")
  }
  if cpu.b != 0 {
    t.Error("b register should be 0")
  }
  if cpu.c != 0 {
    t.Error("c register should be 0")
  }
  if cpu.d != 0 {
    t.Error("d register should be 0")
  }
  if cpu.e != 0 {
    t.Error("e register should be 0")
  }
  if cpu.h != 0 {
    t.Error("h register should be 0")
  }
  if cpu.l != 0 {
    t.Error("l register should be 0")
  }
  if cpu.Pc != 0 {
    t.Error("Pc register should be 0")
  }
  if cpu.sp != 0 {
    t.Error("sp register should be 0")
  }
  if cpu.M != 0 {
    t.Error("M register should be 0")
  }
  if cpu.T != 0 {
    t.Error("T register should be 0")
  }
  if cpu.f != 0 {
    t.Error("f register should be 0")
  }
  if cpu.Stop != false {
    t.Error("Stop flag should be false")
  }
}

func TestWriteIntoBCRegister (t *testing.T) {
  cpu := CPU{}
  cpu.Reset()

  cpu.w_bc(0x0150)

  if cpu.b != 0x01 {
    t.Error("b register should contain value 0x01")
  }
  if cpu.c != 0x50 {
    t.Error("c register should contain value 0x50")
  }
}

func TestNop(t *testing.T) {
  cpu := CPU{}
  cpu.Reset()

  cpu.Dispatch(0x00)

  if cpu.T != 4 {
    t.Error("instruction should consume 1 t cycle")
  }
  if cpu.M != 1 {
    t.Error("instruction should consume 4 m cycle")
  }
  if cpu.Pc != 1 {
    t.Error("instruction should advance Pc by 1")
  }
}

func TestLD(t *testing.T) {
  cpu := CPU{}
  cpu.Reset()

  // LD BC, d16
  cpu.Memory = [65536]byte{ 0x01, 0x50, 0x01 }
  cpu.Dispatch(0x01)

  if cpu.T != 12 {
    t.Error("instruction should consume 12 t cycle")
  }
  if cpu.M != 3 {
    t.Error("instruction should consume 3 m cycle")
  }
  if cpu.b != 0x01 {
    t.Error("b register should contain 0x01")
  }
  if cpu.c != 0x50 {
    t.Error("c register should contain 0x50")
  }

  // LD (BC), A
  cpu.Reset()
  cpu.w_bc(0x0150)
  cpu.a = 0x33
  fmt.Printf("%02x %02x\n", cpu.b, cpu.c)
  cpu.Dispatch(0x02)

  if cpu.T != 8 {
    t.Error("instruction should consume 8 t cycle")
  }
  if cpu.M != 1 {
    t.Error("instruction should consume 1 m cycle")
  }
  if cpu.Memory[0x150] != 0x33 {
    t.Error("there should be the written 0x33 value into the memory address 0x150")
  }
}

func TestInc(t *testing.T) {
  cpu := CPU{}
  cpu.Reset()

  // INC bc
  cpu.w_bc(0x0150)
  cpu.Dispatch(0x03)

  if cpu.T != 8 {
    t.Error("instruction should consume 8 t cycle")
  }
  if cpu.M != 1 {
    t.Error("instruction should consume 1 m cycle")
  }
  if cpu.r_bc() != 0x0151 {
    t.Error("bc register should be increased")
  }
}

