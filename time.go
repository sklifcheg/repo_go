package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.UnixMilli(1710109612937)
	t.Format("2006-01-02 15:04:05")
	fmt.Println(t)
	fmt.Println(t.Format("2006-01-02 15:04:05"))
}
