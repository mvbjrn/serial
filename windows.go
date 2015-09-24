// Copyright 2015 mvbjrn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows !cgo

package serial

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var ()

// TODO flow control

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

	match, _ := regexp.MatchString(".*COM.*", connection.Port)
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
		connection.Port, connection.Baud, connection.DataBit, parity,
		connection.StopBit)
}

// Open a connection.
func (connection *Connection) Open(timeout uint8) error {
	//TODO
}

// Write a byte array to an open connection.
func (connection *Connection) Write(b []byte) (int, error) {
	//TODO
}

// Read from an open connection until the delimiter is reached.
func (connection *Connection) Read(delimiter byte) ([]byte, error) {
	//TODO
}

// Flush the connection, which causes untransmitted or not read
// data to be discarded.
func (connection *Connection) Flush() error {
	//TODO
}

// Close a connection.
func (connection *Connection) Close() error {
	//TODO
}

// functions

// createConnection is the entrence point for the Connection in windows.
func createConnection(port string, baudrate Baud, databit DataBit,
	stopbit StopBit, parity Parity) (*Connection, error) {
	connection := &Connection{Port: port, Baud: baudrate,
		DataBit: databit, StopBit: stopbit, Parity: parity}
	return connection, connection.check()
}
