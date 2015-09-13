[![Build Status](https://travis-ci.org/mvbjrn/serial.svg?branch=master)](https://travis-ci.org/mvbjrn/serial)
[![Coverage](http://gocover.io/_badge/github.com/mvbjrn/serial)](http://gocover.io/github.com/mvbjrn/serial)
[![GoDoc](https://godoc.org/github.com/mvbjrn/serial?status.svg)](https://godoc.org/github.com/mvbjrn/serial)

# serial

This is an extended package for accessing the serial port.
It is inspired by https://github.com/tarm/serial and several forks.

I made this library mainly for educational purposes.
As you can see, I commented my code a lot, just to make myself clear what is happening.

**This package is currently untested.**

## Overview

* Library
* serial2http
* Terminal

## Library

The connection struct is all you need. It encapsulates all parameters necessary
and all functions needed to interact with a serial device. Just initiate a connection and open it.
After that, you are able to read and write from it or to flush the I/O. When you are finished, just close it.

A connection can be initiated by setting all parameters within the code or by loading all parameters from a file.
The library provides several consts to set the baud rate, data and stop bits as well as the parity.

Using the Init function:

`connection, err := serial.InitConnection("/dev/ttyUSB0", serial.Baud115200, serial.DataBit8, serial.StopBit1, serial.ParityNone)`

Loading from a file:

`connection, err := serial.LoadConnection("sample/sample.json")`

The next step is to open the connection:

`err := connection.Open()`

This simply sets up the underlying file and binds the desired parameters to it.

Writing to the serial port is done by just calling the write function:

`n, err := connection.Write([]byte("WriteToPort"))`

This returns the count of transmitted bytes as well as the error.

Reading from the serial port requires a delimiter, which indicates the end of the transmission:

`response, err := connection.Read(10)`

The delimiter is device dependent. Using an ASCII table to find the correct decimal value may help at this point.
The response is a `[]byte`, which contains all data transmitted until the delimiter is reached.

Another way for reading from the port is to use a buffer.

`response, err := connection.ReadToBuffer(256)`

This will read all bytes into a `[]byte` with the given size.

After finishing reading and writing, the connection can be flushed and closed:

`err := connection.Flush()`

This causes the I/O to discard untransmitted and unread data.

`err := connection.Close()`

This causes the connection to be closed.

### Possible Errors

* invalid baud rate
* invalid data bits
* invalid stop bits
* invalid parity
* connection not open
  * should only occur when trying to write or read
  * solution: execute connection.Open() before trying to write or read

## serial2http

## Terminal
