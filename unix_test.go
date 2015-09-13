// licence goes here

package serial

import (
	"testing"
	"time"
)

// TODO: documentation
func TestConnection(t *testing.T) {
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

	ch := make(chan int, 1)
	go func() {
		buf := make([]byte, 128)
		var readCount int
		for {
			n, err := connection2.ReadToBuffer(buf)
			if err != nil {
				t.Fatal(err)
			}
			readCount++
			t.Logf("Read %v %v bytes: % 02x %s", readCount, n, buf[:n], buf[:n])
			select {
			case <-ch:
				ch <- readCount
				close(ch)
			default:
			}
		}
	}()

	if _, err = connection1.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if _, err = connection1.Write([]byte(" ")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if _, err = connection1.Write([]byte("world")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second / 10)

	ch <- 0
	connection1.Write([]byte(" ")) // We could be blocked in the read without this
	c := <-ch
	exp := 5
	if c >= exp {
		t.Fatalf("Expected less than %v read, got %v", exp, c)
	}

	connection1.Flush()
	connection1.Close()

	connection2.Flush()
	connection2.Close()
}
