// licence goes here

package serial

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// const

// Baud is the unit for the symbol rate. It describes the number of symbols transmitted per second.
type Baud uint32

const (
	// Baud4800 defines a transmission rate of 4800 symbols per second.
	Baud4800 = 4800
	// Baud9600 defines a transmission rate of 9600 symbols per second.
	Baud9600 = 9600
	// Baud19200 defines a transmission rate of 19200 symbols per second.
	Baud19200 = 19200
	// Baud38400 defines a transmission rate of 38400 symbols per second.
	Baud38400 = 38400
	// Baud57600 defines a transmission rate of 57600 symbols per second.
	Baud57600 = 57600
	// Baud115200 defines a transmission rate of 115200 symbols per second.
	Baud115200 = 115200
)

// DataBit is the number of bits representing a character.
type DataBit byte

const (
	// DataBit5 stands for a character length of five bits.
	DataBit5 = DataBit(iota + 5)
	// DataBit6 stands for a character length of six bits.
	DataBit6
	// DataBit7 stands for a character length of seven bits.
	DataBit7
	// DataBit8 stands for a character length of eight bits.
	DataBit8
	// DataBit9 stands for a character length of nine bits.
	DataBit9
)

// StopBit is the number of bits being send at the end of every character.
type StopBit byte

const (
	// StopBit1 represents a single bit being send as stopbit.
	StopBit1 = StopBit(iota + 1)
	// StopBit2 represents two bits being send as stopbit.
	StopBit2
)

// Parity is the method for detecting transmission errors.
type Parity byte

const (
	// ParityNone indicates that no error detection is being used.
	ParityNone = Parity(iota)
	// ParityEven indicates that a bit is added to even out the bit count.
	ParityEven
	// ParityOdd indicates that a bit is added to provide an odd bit count.
	ParityOdd
)

// TODO flow control

// var
var (
	errBaud    = errors.New("serial configuration error: bad baud rate (4800, 9600, 19200, 38400, 57600, 115200)")
	errDataBit = errors.New("serial configuration error: bad number of data bits (5, 6, 7, 8, 9)")
	errStopBit = errors.New("serial configuration error: bad number of stop bits (1, 2)")
	errParity  = errors.New("serial configuration error: bad parity (0 - None, 1 - Even, 2 - Odd)")
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
}

func (connection *Connection) check() error {

	switch connection.Baud {
	case Baud115200, Baud57600, Baud38400, Baud19200, Baud9600, Baud4800:
	default:
		return errBaud
	}

	switch connection.DataBit {
	case DataBit5, DataBit6, DataBit7, DataBit8, DataBit9:
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

// functions

// Init provides a connection with the given parameters.
func Init(port string, baudrate Baud, databit DataBit, stopbit StopBit, parity Parity) (*Connection, error) {
	connection := &Connection{Port: port, Baud: baudrate, DataBit: databit, StopBit: stopbit, Parity: parity}
	return connection, connection.check()
}

// LoadFile provides a connection with the parameters being loaded from a json file.
func LoadFile(path string) (*Connection, error) {
	var connection *Connection

	file, err := os.Open(path)
	if err != nil {
		return connection, err
	}

	dec := json.NewDecoder(file)
	for {

		if err := dec.Decode(&connection); err == io.EOF {
			break
		} else if err != nil {
			return connection, err
		}
	}

	return connection, connection.check()
}
