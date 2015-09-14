// Copyright 2015 mvbjrn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serial

import "testing"

// TestConnection is part of a loopback test. Additional information is provided in the repo-wiki.
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
	err = connection1.Open(0)
	if err != nil {
		t.Fatal(err)
	}

	err = connection2.Open(0)
	if err != nil {
		t.Fatal(err)
	}

	//TODO write the actual test. One connection should write, the other one should read. then compare the results!

	connection1.Flush()
	connection1.Close()

	connection2.Flush()
	connection2.Close()
}
