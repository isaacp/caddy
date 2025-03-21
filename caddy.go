package caddy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	flag struct {
		value string
	}

	argument struct {
		value string
	}

	Command struct {
		name string
	}

	Output struct {
		Error error
		Value string
	}
)

func (arg *argument) expand() string {
	return os.Expand(arg.value, func(arg string) string {
		if value, ok := os.LookupEnv(arg); ok {
			return value
		}

		return arg
	})
}

func Execute[T any](inputs ...string) Output {
	for index, inp := range inputs {
		if strings.HasPrefix(inp, "$") {
			rg := argument{value: inp}
			inputs[index] = rg.expand()
		}
	}
	cmd := exec.Command(inputs[0], inputs[1:]...)
	bytes, err := cmd.Output()
	if err != nil {
		return Output{
			Value: "",
			Error: err,
		}
	}
	return Output{
		Value: string(bytes),
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("version 0.1")
		return
	}

	os.Args = os.Args[1:]
	flags := os.Args[0]
	repeat := false

	if strings.HasPrefix(flags, "-") {
		if len(os.Args) == 1 {
			log.Println("no command defined")
			return
		}
		repeat = strings.Contains(flags, "f")
		os.Args = os.Args[1:]
	}

	for {
		output := Execute[string](os.Args...)
		log.Printf("%+v", output.Value)

		if !repeat {
			break
		}

		time.Sleep(time.Second)
	}
}
