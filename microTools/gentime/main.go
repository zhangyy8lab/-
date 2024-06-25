package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Current UTC Timestamp:", time.Now().UTC().Unix())
	fmt.Println("Current  +8 Timestamp:", time.Now().Add(8*time.Hour).Unix())
}
