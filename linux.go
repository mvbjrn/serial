// licence goes here

// +build linux,!cgo

package serial

import (
	"bufio"
	"errors"
	"os"
	"syscall"
)

var (
	errConnOpen = errors.New("serial connection error: connection is not open")
)

// structs and its functions

// Open a connection.
func (connection *Connection) Open() error {

	// syscall: https://golang.org/pkg/syscall/

	// The serial port is basically a file we are writing to and reading from.
	// 	O_RDWR allows us to read and write.
	// 	O_NOCTTY prevents the device from controlling the terminal.
	// 	O_NONBLOCK prevents the system from blocking for a long time.
	connection.f, err = os.OpenFile(connection.Port, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return err
	}

	connection.open = true
	return nil
}

// Write a byte array to an open connection.
func (connection *Connection) Write(b []byte) (int, error) {
	if connection.open {
		return connection.f.Write(b)
	}

	return _, errConnOpen
}

// Read from an open connection until the delimter is reached.
func (connection *Connection) Read(delimter byte) ([]byte, error) {
	if connection.open {
		reader := bufio.NewReader(connection.f)
		return reader.ReadBytes(delimiter)
	}

	return _, errConnOpen
}

// Close a connection.
func (connection *Connection) Close() error {
	err := connection.f.Close()
	connection.open = false

	return err
}
