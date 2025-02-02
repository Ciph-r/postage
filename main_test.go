package main

import "testing"

func TestFoo(t *testing.T) {
	t.Fatal("this should block the pr.")
}
