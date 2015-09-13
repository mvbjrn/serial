// licence goes here

package serial

import "testing"

// TODO: documentation
func TestConnection(t *testing.T) {
	t.Skip()
	connection1, err := InitConnection("/dev/ttyUSB0", Baud115200, DataBit8, StopBit1, ParityNone)
	if err != nil {
		t.Fatal(err)
	}

	connection2, err := InitConnection("/dev/ttyUSB1", Baud115200, DataBit8, StopBit1, ParityNone)
	if err != nil {
		t.Fatal(err)
	}

	//s1, err := OpenPort(c0)
	err = connection1.Open()
	if err != nil {
		t.Fatal(err)
	}

	err = connection2.Open()
	if err != nil {
		t.Fatal(err)
	}

	//TODO write the actual test. One connection should write, the other one should read. then compare the results!

	connection1.Flush()
	connection1.Close()

	connection2.Flush()
	connection2.Close()
}
