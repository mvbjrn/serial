// licence goes here

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
	t.Skip("skipping test.")
	connection, err := Init("/dev/ttyUSB0", Baud115200, DataBit8, StopBit1, ParityNone)
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
	t.Skip("skipping test.")
	connection1, err := Init("/dev/ttyUSB0", Baud115200, DataBit8, StopBit1, ParityNone)
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

	connection2, err := Load(sample)
	if err != nil {
		t.Error(err)
	}

	if connection1.String() != connection2.String() {
		t.Error("connection1 defers from connection2")
	}
}
