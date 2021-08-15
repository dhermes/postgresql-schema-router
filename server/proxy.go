package server

import (
	"io"
	"net"
	"sync"
	"time"
)

const (
	readTimeout = 250 * time.Millisecond
)

type forwardState struct {
	Mutex  sync.RWMutex
	Errors []error
	Done   bool
}

func (fs *forwardState) AddError(err error) {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()
	fs.Errors = append(fs.Errors, err)
	fs.Done = true
}

func (fs *forwardState) IsDone() bool {
	fs.Mutex.RLock()
	defer fs.Mutex.RUnlock()
	return fs.Done
}

func (fs *forwardState) MarkDone() {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()
	fs.Done = true
}

func forward(wg *sync.WaitGroup, r, w *net.TCPConn, fs *forwardState) {
	defer wg.Done()

	data := make([]byte, 4096)
	for {
		if fs.IsDone() {
			return
		}

		readDeadline := time.Now().Add(readTimeout)
		err := r.SetReadDeadline(readDeadline)
		if err != nil {
			fs.AddError(err)
			return
		}

		n, err := r.Read(data)
		if err == io.EOF {
			fs.MarkDone()
			return
		}
		if err != nil {
			if isTimeout(err) {
				continue
			}
			fs.AddError(err)
			return
		}

		_, err = w.Write(data[:n])
		if err != nil {
			fs.AddError(err)
			return
		}
	}
}

// proxyInternal is the underlying implementation for `proxy()`, but
// it does not have to do any extra resolution of errors.
func proxyInternal(tc *net.TCPConn, c Config) (err error) {
	var sc *net.TCPConn
	defer func() {
		if sc == nil {
			return
		}
		err = appendErrs(err, sc.Close())
	}()

	addr, err := net.ResolveTCPAddr("tcp", c.RemoteAddr)
	if err != nil {
		return
	}

	sc, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	fs := forwardState{}
	go forward(&wg, tc, sc, &fs)
	go forward(&wg, sc, tc, &fs)
	wg.Wait()

	err = appendErrs(fs.Errors...)
	return nil
}

// proxy is the "pristine" function to be directly used in a `goroutine`.
// It is fully responsible for cleaning up after itself.
func proxy(tc *net.TCPConn, c Config) {
	err := proxyInternal(tc, c)
	if err == nil {
		return
	}
	// LOG-TODO: Do something with the error
}
