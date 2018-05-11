package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	c *exec.Cmd

	rootCmd = &cobra.Command{
		Use:   "jlp",
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

			green := color.New(color.FgGreen).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			_ = red

			go func() {
				outscan := bufio.NewScanner(outpipe)
				for {
					chk := outscan.Scan()
					if !chk {
						if outscan.Err() != nil {
							log.Fatalln(outscan.Err())
						}

						break
					}

					stxt := outscan.Text()
					if stxt[0] == '{' || stxt[0] == '[' {
						stxt = pretty(stxt)
					}

					log.Println(green("[stdout]"), stxt)
				}
			}()

			go func() {
				errscan := bufio.NewScanner(errpipe)
				for {
					chk := errscan.Scan()
					if !chk {
						if errscan.Err() != nil {
							log.Fatalln(errscan.Err())
						}

						break
					}

					stxt := errscan.Text()
					if stxt[0] == '{' || stxt[0] == '[' {
						stxt = pretty(stxt)
					}

					log.Println(red("[stderr]"), stxt)
				}
			}()

			c.Wait()
		},
	}
)

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
		return err.Error()
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

	rootCmd.Execute()
}
