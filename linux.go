// licence goes here

// +build linux,!cgo

package serial

import (
	"bufio"
	"errors"
	"os"
	"syscall"
	"unsafe"
)

var (
	errConnOpen = errors.New("serial connection error: connection is not open")
	baudrates   = map[Baud]uint32{
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

// Open a connection.
func (connection *Connection) Open() error {

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
		if err != nil && f != nil {
			f.Close()
			return err
		}
	}()

	// Create a file descriptor.
	fd := connection.f.Fd()

	// Create a termios with the given parameters.
	t := syscall.Termios{}

	// Setup the baud rate.
	baudrate, err := baudrates[connection.Baud]
	if err != nil {
		return errBaud
	}

	t.Cflag |= baudrate
	t.Ispeed = baudrate
	t.Ospeed = baudrate

	// Setup stop bits.
	switch connection.StopBit {
	case StopBits1:
		t.Cflag &^= syscall.CSTOPB // CSTOPB = 0x40
	case StopBits2:
		t.Cflag |= syscall.CSTOPB
	default:
		return errStopBit
	}

	// Setup data bits.
	databit, err := databits[connection.DataBit]
	if err != nil {
		return errDataBit
	}
	t.Cflag |= databit

	// Setup parity.
	switch connection.Parity {
	case ParityNone:
		t.Cflag &^= syscall.PARENB // PARENB = 0x100
	case ParityEven:
		t.Cflag |= syscall.PARENB
	case ParityOdd:
		t.Cflag |= syscall.PARENB
		t.Cflag |= syscall.PARODD // PARODD = 0x200
	default:
		return errParity
	}

	// Execute TCSETS with the modified termios structure to apply changes.
	if _, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	); errno != 0 {
		return nil, errno
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
