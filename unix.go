// licence goes here

// +build linux, !cgo

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
		syscall.SYS_IOCTL,                 // device-specific input/output operations
		uintptr(connection.f.Fd()),        // open file descriptor
		uintptr(syscall.TCSETS),           // a request code number to set the current serial port settings
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
