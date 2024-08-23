package tests

import "testing"

const ret string = "Helo"

func TestAssertHello(t *testing.T) {
	if ret != "Hello" {
		t.Fatalf("should have been test")
	}
}

//enter on this folder and type on terminal "go test"
