package main

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
)

// NOTE: Ensure that
//       - `B` satisfies `io.Reader`
//       - `B` satisfies `io.Writer`
var (
	_ io.Reader = (*B)(nil)
	_ io.Writer = (*B)(nil)
)

type B []byte

func NewB(s []byte) B {
	return B(s)
}

func (b *B) Read(p []byte) (int, error) {
	s := *b
	m := intMin(len(p), len(s))
	if m == 0 {
		return 0, nil
	}
	copy(p[:m], s[:m])

	// Free the memory in `b` that was just consumed.
	newS := make([]byte, len(s)-m)
	copy(newS, s[m:])
	*b = newS

	return m, nil
}

func (b *B) Discard(n int) error {
	s := *b
	if n > len(s) {
		return bufio.ErrBufferFull
	}
	if n < 0 {
		return bufio.ErrNegativeCount
	}

	// Free the memory in `b` that was just consumed.
	newS := make([]byte, len(s)-n)
	copy(newS, s[n:])
	*b = newS
	return nil
}

func (b *B) Peek(n int) ([]byte, error) {
	s := *b
	if n > len(s) {
		return s, bufio.ErrBufferFull
	}
	if n < 0 {
		return nil, bufio.ErrNegativeCount
	}
	return s[:n], nil
}

func (b *B) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *B) Len() int {
	return len(*b)
}

func (b *B) Cap() int {
	return cap(*b)
}

func main() {
	var a1, h1, t1, a2, h2, t2 uint64
	ms := runtime.MemStats{}
	var b B

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a1 = ms.Alloc
	h1 = ms.HeapAlloc
	t1 = ms.TotalAlloc
	fmt.Printf("0: Alloc = %d; HeapAlloc = %d; TotalAlloc = %d\n", a1, h1, t1)

	runtime.GC()
	b = NewB(make([]byte, 0x10000))
	fmt.Printf("len(b) = %d; cap(b) = %d\n", b.Len(), b.Cap())
	runtime.ReadMemStats(&ms)
	a1 = ms.Alloc
	h1 = ms.HeapAlloc
	t1 = ms.TotalAlloc
	fmt.Printf("1: \u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", -diff(a2, a1), -diff(h2, h1), -diff(t2, t1))

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a2 = ms.Alloc
	h2 = ms.HeapAlloc
	t2 = ms.TotalAlloc
	fmt.Printf("2: \u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", diff(a2, a1), diff(h2, h1), diff(t2, t1))
	b.Discard(0xc000)
	fmt.Printf("len(b) = %d; cap(b) = %d\n", b.Len(), b.Cap())

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a1 = ms.Alloc
	h1 = ms.HeapAlloc
	t1 = ms.TotalAlloc
	fmt.Printf("3: \u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", -diff(a2, a1), -diff(h2, h1), -diff(t2, t1))
	b.Discard(0x3f00)
	fmt.Printf("len(b) = %d; cap(b) = %d\n", b.Len(), b.Cap())

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a2 = ms.Alloc
	h2 = ms.HeapAlloc
	t2 = ms.TotalAlloc
	fmt.Printf("4: \u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", diff(a2, a1), diff(h2, h1), diff(t2, t1))
	b.Write([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse ci..."))
	fmt.Printf("len(b) = %d; cap(b) = %d\n", b.Len(), b.Cap())

	runtime.GC()
	runtime.ReadMemStats(&ms)
	a1 = ms.Alloc
	h1 = ms.HeapAlloc
	t1 = ms.TotalAlloc
	fmt.Printf("5: \u0394 Alloc = %+d; \u0394 HeapAlloc = %+d; \u0394 TotalAlloc = %+d\n", -diff(a2, a1), -diff(h2, h1), -diff(t2, t1))
}

func diff(u1, u2 uint64) int64 {
	if u1 > u2 {
		return int64(u1 - u2)
	}
	return -int64(u2 - u1)
}

func intMin(i1, i2 int) int {
	if i1 < i2 {
		return i1
	}
	return i2
}
