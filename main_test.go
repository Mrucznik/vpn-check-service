package main

import (
	"fmt"
	"testing"
)

func TestGetIps(t *testing.T) {
	for _, ip := range getIps() {
		fmt.Println(ip)
	}
}

