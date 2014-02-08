package main

import (
	// #include <SDL2/SDL.h>
	// #cgo LDFLAGS: -L/usr/local/lib -lSDL2 -Wl,-rpath=/usr/local/lib
	"C"
	"fmt"
	"os"
	. "gogb/z80"
)

func main() {
	z80 := Z80{Stop: false}
	z80.Init()
	
	// mmu.loadRom(os.Args[1])
	z80.LoadBios()

	codeLength := len(z80.Memory)
	fmt.Printf("Code is %d bytes large\n", codeLength)

	z80.Reset()

	for int(z80.Pc) < codeLength && !z80.Stop {
		op := z80.Memory[z80.Pc]

		z80.DumpOp(op)

		z80.Dispatch(op)
	}

	n, _ := C.SDL_Init(C.SDL_INIT_EVERYTHING)
	if n != 0 {
		fmt.Println("Error initializing SDL: %s", C.SDL_GetError())
	}
	defer C.SDL_Quit()

	win := C.SDL_CreateWindow(C.CString("Hello world!"), 100, 100, 640, 480, C.SDL_WINDOW_SHOWN)
	if win == nil {
		fmt.Println("Error creating a window: %s", C.SDL_GetError())
	}		
	
	os.Exit(0)
}
