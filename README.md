# serial

This is an extended package for accessing the serial port.
It is inspired by https://github.com/tarm/serial and several forks.

I made this library mainly for educational purposes.
As you can see, I commented my code alot, just to make myself clear what is happening.

## Overview

* Library
* serial2http
* Terminal

## Library

The connection struct is all you need. It encapsulates all paramaters necessary
and all functions needed to interact with a serial device. Just initiate a connection and open it.
After that your are able to read and write from it. When you are finished, just close it.

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
