package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	go func() {
		f, err := os.Open("test.txt")
		if err != nil {
			fmt.Println("Error opening test.txt (call 1)")
			return
		}
		defer f.Close()
		_, err = f.Seek(6, 0)
		if err != nil {
			fmt.Println("Error seeking in test.txt (call 1)")
			return
		}
		out, err := os.Create("out1.txt")
		if err != nil {
			return
		}
		defer out.Close()
		io.CopyN(out, f, 5)
	}()
	go func() {
		f, err := os.Open("test.txt")
		if err != nil {
			fmt.Println("Error opening test.txt")
			return
		}
		defer f.Close()
		_, err = f.Seek(11, 0)
		if err != nil {
			fmt.Println("Error seeking in test.txt (call 2)")
			return
		}
		out, err := os.Create("out2.txt")
		if err != nil {
			return
		}
		defer out.Close()
		io.CopyN(out, f, 4)
	}()
	time.Sleep(2 * time.Second)
}

/*
import (
	"fmt"
	"os"

	"github.com/sscheele/GoDistribute/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
*/
