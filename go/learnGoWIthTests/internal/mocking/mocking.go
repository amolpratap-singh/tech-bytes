package mocking

import (
	"fmt"
	"io"
	"time"
)


const finalWord = "Go!"
const countdownStart = 3

type Sleeper interface {
	Sleep()
}

func CountDown(out io.Writer) {

	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		time.Sleep(1 * time.Second)
	}
	fmt.Fprint(out, finalWord)

}