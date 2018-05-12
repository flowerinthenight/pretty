package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type fnSprintFunc func(a ...interface{}) string // color's SprintFunc's signature

var (
	c *exec.Cmd

	stern bool
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()

	colors = []color.Attribute{
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
		color.FgCyan,
		color.FgWhite,
		color.FgRed,
		color.FgGreen,
	}

	rootCmd = &cobra.Command{
		Use:   "pretty",
		Short: "JSON log prettifier wrapper tool",
		Long:  "JSON log prettifier wrapper tool.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				c = exec.Command(args[0])
			}

			if len(args) > 1 {
				c = exec.Command(args[0], args[1:]...)
			}

			if c == nil {
				log.Fatalln("no args")
			}

			outpipe, err := c.StdoutPipe()
			if err != nil {
				log.Fatalln(err)
			}

			errpipe, err := c.StderrPipe()
			if err != nil {
				log.Fatalln(err)
			}

			err = c.Start()
			if err != nil {
				log.Fatalln(err)
			}

			go func() {
				outscan := bufio.NewScanner(outpipe)
				print(green("[stdout]"), outscan)
			}()

			go func() {
				errscan := bufio.NewScanner(errpipe)
				print(red("[stderr]"), errscan)
			}()

			c.Wait()
		},
	}
)

func print(outpre string, scan *bufio.Scanner) {
	var ci int
	var cm map[string]fnSprintFunc

	if stern {
		cm = make(map[string]fnSprintFunc)
	}

	for {
		chk := scan.Scan()
		if !chk {
			if scan.Err() != nil {
				log.Fatalln(scan.Err())
			}
		}

		b := scan.Bytes()
		pre, s := prepare(b)
		if pre != "" {
			if stern {
				fnclr, ok := cm[pre]
				if ok {
					pre = fnclr(pre)
				} else {
					fn := color.New(colors[ci]).SprintFunc()
					cm[pre] = fn
					pre = fn(pre)
					ci += 1
					if ci >= len(colors) {
						ci = 0
					}
				}
			}

			log.Println(outpre, pre, s)
		} else {
			log.Println(outpre, s)
		}
	}
}

// Returns prefix (before the JSON part, if any), and the JSON string.
func prepare(b []byte) (string, string) {
	var prefix string

	s := string(b)
	i1 := strings.Index(s, "[")
	i2 := strings.Index(s, "{")
	if i1 < 0 && i2 < 0 {
		return prefix, s
	}

	if (i1 * i2) == 0 {
		return prefix, pretty(s)
	}

	i := i1
	if (i1 * i2) < 0 {
		if i2 > i {
			i = i2
		}
	} else {
		if i2 < i {
			i = i2
		}
	}

	if i > 0 {
		prefix = s[0:i]
	}

	return prefix, pretty(s[i:])
}

func pretty(v interface{}) string {
	var out bytes.Buffer
	var b []byte

	_, ok := v.(string)
	if !ok {
		tmp, err := json.Marshal(v)
		if err != nil {
			return err.Error()
		}

		b = tmp
	} else {
		b = []byte(v.(string))
	}

	err := json.Indent(&out, b, "", "  ")
	if err != nil {
		return v.(string)
	}

	return out.String()
}

func main() {
	go func() {
		s := make(chan os.Signal)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		sig := errors.Errorf("%s", <-s)
		_ = sig

		if c != nil {
			err := c.Process.Signal(syscall.SIGTERM)
			if err != nil {
				log.Println("failed to terminate process, force kill...")
				_ = c.Process.Signal(syscall.SIGKILL)
			}
		}

		os.Exit(0)
	}()

	rootCmd.Flags().BoolVar(&stern, "stern", stern, "prefix color if using stern")
	rootCmd.Execute()
}
