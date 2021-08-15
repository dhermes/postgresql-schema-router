package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/dhermes/postgresql-schema-router/postgres"
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

type packetInspector func(data []byte)

func forward(wg *sync.WaitGroup, r, w *net.TCPConn, fs *forwardState, pi packetInspector) {
	defer wg.Done()

	data := make([]byte, 65536)
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
		// Ensure we have read a "complete" TCP packet, with a limit on the size.
		if n >= 65536 {
			fs.AddError(fmt.Errorf("%w, exceeds 65536 bytes", ErrPacketTooLarge))
			return
		}

		if pi != nil {
			pi(data[:n])
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
	go forward(&wg, tc, sc, &fs, inspectPG) // Client->Proxy->Remote
	go forward(&wg, sc, tc, &fs, nil)       // Remote->Proxy->Client
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

func inspectPG(chunk []byte) {
	fm, err := postgres.ParseChunk(chunk)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Failed to parse TCP chunk as PostgreSQL frontend message; %v",
			err,
		)
		return
	}

	fmt.Printf("FrontendMessage: %#v\n", fm)
}
