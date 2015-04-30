package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func main() {
	cmd := flag.String("cmd", "", "Command excuted")
	eT := flag.Int64("e", 0, "Expect Time Process Run")
	flag.Parse()
	var t0 time.Time
	var d0 int64

	// Check if command and expect time is Zero
	if *cmd == "" || *eT < 1 {
		err := "Command or Expect Time is nil"
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmds := strings.Split(*cmd, " ")

	// Function that excute command with Lock
	exe_cmd := func(wg *sync.WaitGroup, cmds []string) {
		// fmt.Println(cmd)
		go func() {
			t0 = time.Now()
		}()
		cmd := cmds[0]
		args := []string{}
		if len(cmds) > 1 {
			args = cmds[1:]
		}

		_, err := exec.Command(cmd, args...).Output()
		if err != nil {
			fmt.Println("error occured")
			fmt.Printf("%s", err)
		}
		// fmt.Printf("%s", out)
		// check RunTime in end of function
		go func() {
			for {
				if d0 = int64(time.Now().Sub(t0).Nanoseconds() / 1000000); d0 < *eT*1000 {
					err := "Command doesnt run enough time"
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(1)
				}
			}
		}()
		wg.Done()
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go exe_cmd(wg, cmds)
	wg.Wait()

	fmt.Printf("%v\n", d0)

}
