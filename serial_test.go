// Copyright 2015 mvbjrn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serial

import (
	"os"
	"testing"
)

// const
const ()

// var
var (
	sample = "sample.json"
)

// structs and its functions

// functions

func TestSave(t *testing.T) {
	connection, err := InitConnection("/dev/ttyUSB0", Baud115200, DataBit8, StopBit1, ParityNone)
	if err != nil {
		t.Error(err)
	}

	err = connection.Save(sample)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = os.Remove(sample)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestLoad(t *testing.T) {
	connection1, err := InitConnection("/dev/ttyUSB0", Baud115200, DataBit8, StopBit1, ParityNone)
	if err != nil {
		t.Error(err)
	}

	err = connection1.Save(sample)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = os.Remove(sample)
		if err != nil {
			t.Error(err)
		}
	}()

	connection2, err := LoadConnection(sample)
	if err != nil {
		t.Error(err)
	}

	if connection1.String() != connection2.String() {
		t.Error("connection1 defers from connection2")
	}
}
