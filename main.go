package main

import (
	"bufio"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "jlp",
		Short: "JSON log prettifier wrapper tool",
		Long:  "JSON log prettifier wrapper tool.",
		Run: func(cmd *cobra.Command, args []string) {
			var c *exec.Cmd

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
				for {
					chk := outscan.Scan()
					if !chk {
						if outscan.Err() != nil {
							log.Fatalln(outscan.Err())
						}

						break
					}

					stxt := outscan.Text()
					log.Println(stxt)
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
					log.Println(stxt)
				}
			}()

			c.Wait()
		},
	}
)

func main() {
	rootCmd.Execute()
}
