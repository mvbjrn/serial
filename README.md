# serial

An extended package for accessing the serial port.
It is inspired by https://github.com/tarm/serial.

## Overview

* Library
* serial2http
* Terminal

## Library

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
