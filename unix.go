// Copyright 2015 mvbjrn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows !cgo
// +build 386 amd64 arm

package serial

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"syscall"
	"unsafe"
)

var (
	baudrates = map[Baud]uint32{
		Baud4800:   syscall.B4800,
		Baud9600:   syscall.B9600,
		Baud19200:  syscall.B19200,
		Baud38400:  syscall.B38400,
		Baud57600:  syscall.B57600,
		Baud115200: syscall.B115200,
	}
	databits = map[DataBit]uint32{
		DataBit5: syscall.CS5,
		DataBit6: syscall.CS6,
		DataBit7: syscall.CS7,
		DataBit8: syscall.CS8,
	}
)

// structs and its functions

// Connection represents a serial connection with all parameters.
type Connection struct {
	Port    string
	Baud    Baud
	DataBit DataBit
	StopBit StopBit
	Parity  Parity
	f       *os.File
	isOpen  bool
}

func (connection *Connection) check() error {

	match, _ := regexp.MatchString("/dev/tty.*", connection.Port)
	if !match {
		return errPort
	}

	switch connection.Baud {
	case Baud115200, Baud57600, Baud38400, Baud19200, Baud9600, Baud4800:
	default:
		return errBaud
	}

	switch connection.DataBit {
	case DataBit5, DataBit6, DataBit7, DataBit8:
	default:
		return errDataBit
	}

	switch connection.StopBit {
	case StopBit1, StopBit2:
	default:
		return errStopBit
	}

	switch connection.Parity {
	case ParityNone, ParityEven, ParityOdd:
	default:
		return errParity
	}

	return nil
}

// Save a connection to a json file.
func (connection *Connection) Save(path string) error {
	json, err := json.Marshal(connection)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, json, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (connection *Connection) String() string {

	var parity string
	switch connection.Parity {
	case ParityNone:
		parity = "N"
	case ParityEven:
		parity = "E"
	case ParityOdd:
		parity = "O"
	}

	return fmt.Sprintf("port: %s, baud rate:%d, parameters: %d%s%d",
		connection.Port, connection.Baud, connection.DataBit, parity, connection.StopBit)
}

// Open a connection.
func (connection *Connection) Open() error {

	var err error

	// The serial port is basically a file we are writing to and reading from.
	// 	O_RDWR allows the program to read and write the file.
	// 	O_NOCTTY prevents the device from controlling the terminal.
	// 	O_NONBLOCK prevents the system from blocking for a long time.
	connection.f, err = os.OpenFile(connection.Port, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return err
	}

	// Close the file on error occurrence.
	defer func() {
		if err != nil && connection.f != nil {
			connection.f.Close()
		}
	}()

	// Create a plain termios, which allows the program to execute input/output operations.
	termios := syscall.Termios{}

	// Setup the baud rate in the termios structure.
	baudrate := baudrates[connection.Baud]

	termios.Cflag |= baudrate
	termios.Ispeed = baudrate
	termios.Ospeed = baudrate

	// Setup stop bits in the termios structure.
	switch connection.StopBit {
	case StopBit1:
		termios.Cflag &^= syscall.CSTOPB // CSTOPB = 0x40
	case StopBit2:
		termios.Cflag |= syscall.CSTOPB
	default:
		return errStopBit
	}

	// Setup data bits in the termios structure.
	databit := databits[connection.DataBit]
	termios.Cflag |= databit

	// Setup parity in the termios structure.
	switch connection.Parity {
	case ParityNone:
		termios.Cflag &^= syscall.PARENB // PARENB = 0x100
	case ParityEven:
		termios.Cflag |= syscall.PARENB
	case ParityOdd:
		termios.Cflag |= syscall.PARENB
		termios.Cflag |= syscall.PARODD // PARODD = 0x200
	default:
		return errParity
	}

	// Execute IOCTL with the modified termios structure to apply the changes.
	if _, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,          // device-specific input/output operations
		uintptr(connection.f.Fd()), // open file descriptor
		uintptr(syscall.TCSETS),    // a request code number to set the current serial port settings
		//TODO: it looks like syscall.TCSETS is not available under freebsd and darwin. Is this a bug?
		uintptr(unsafe.Pointer(&termios)), // a pointer to the termios structure
		0,
		0,
		0,
	); errno != 0 {
		return errno
	}

	connection.isOpen = true
	return nil
}

// Write a byte array to an open connection.
func (connection *Connection) Write(b []byte) (int, error) {
	if connection.isOpen {
		return connection.f.Write(b)
	}

	return 0, errConnOpen
}

// Read from an open connection until the delimter is reached.
func (connection *Connection) Read(delimiter byte) ([]byte, error) {
	if connection.isOpen {
		reader := bufio.NewReader(connection.f)
		return reader.ReadBytes(delimiter)
	}

	return nil, errConnOpen
}

// ReadToBuffer reads from an open connection into a []byte buffer with the given size.
func (connection *Connection) ReadToBuffer(size int) ([]byte, error) {
	buffer := make([]byte, size)
	//TODO do something with bytes read
	_, err := connection.f.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// Flush the connection, which causes untransmitted or not read data to be discarded.
func (connection *Connection) Flush() error {
	if connection.isOpen {
		_, _, err := syscall.Syscall(
			syscall.SYS_IOCTL,          // device-specific input/output operations
			uintptr(connection.f.Fd()), // open file descriptor
			uintptr(syscall.TCIOFLUSH), // a request code number to flush input/output
			uintptr(0),                 // a pointer to data, not needed here
		)
		return err
	}
	return errConnOpen
}

// Close a connection.
func (connection *Connection) Close() error {
	err := connection.f.Close()
	connection.isOpen = false

	return err
}

func createConnection(port string, baudrate Baud, databit DataBit, stopbit StopBit, parity Parity) (*Connection, error) {
	connection := &Connection{Port: port, Baud: baudrate, DataBit: databit, StopBit: stopbit, Parity: parity}
	return connection, connection.check()
}
