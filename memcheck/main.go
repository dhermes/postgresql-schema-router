package main

import (
	"fmt"
	"runtime"
)

func main() {
	var a1, h1, t1, a2, h2, t2 uint64
	ms := runtime.MemStats{}
	var b []byte

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a1 = ms.Alloc
	h1 = ms.HeapAlloc
	t1 = ms.TotalAlloc
	fmt.Printf("Alloc = %d; HeapAlloc = %d; TotalAlloc = %d\n", a1, h1, t1)

	runtime.GC()
	b = make([]byte, 0x10000)
	fmt.Printf("len(b) = %d; cap(b) = %d\n", len(b), cap(b))
	runtime.ReadMemStats(&ms)
	a1 = ms.Alloc
	h1 = ms.HeapAlloc
	t1 = ms.TotalAlloc
	fmt.Printf("\u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", -diff(a2, a1), -diff(h2, h1), -diff(t2, t1))

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a2 = ms.Alloc
	h2 = ms.HeapAlloc
	t2 = ms.TotalAlloc
	fmt.Printf("\u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", diff(a2, a1), diff(h2, h1), diff(t2, t1))
	bMove := make([]byte, len(b)-0xc000)
	copy(bMove, b[0xc000:])
	b = bMove
	// NOTE: b = b[0xc000:] will not free
	fmt.Printf("len(b) = %d; cap(b) = %d\n", len(b), cap(b))

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a1 = ms.Alloc
	h1 = ms.HeapAlloc
	t1 = ms.TotalAlloc
	fmt.Printf("\u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", -diff(a2, a1), -diff(h2, h1), -diff(t2, t1))
}

func diff(u1, u2 uint64) int64 {
	if u1 > u2 {
		return int64(u1 - u2)
	}
	return -int64(u2 - u1)
}
