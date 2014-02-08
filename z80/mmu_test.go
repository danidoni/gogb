package cpu

import (
  "testing"
)

func TestReadWordFromMMU(t *testing.T) {
  mmu := MMU{}
  mmu.Memory = [65536]byte{ 0x50, 0x01 }

  if mmu.rw(0x00) != 0x0150 {
    t.Error("word at 0x00 should be 0x0150")
  }
}
