// Copyright © 2019 Jonathan Fontes <jonathan.fontes@creativecodesolutions.pt>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

// TODO
// https://github.com/op/go-logging

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/fatih/color"
	"github.com/mattn/go-zglob"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var create bool

type fileDigest struct {
	file   string
	digest string
}

var list []fileDigest
var listCheck []fileDigest

// digestCmd represents the digest command
var digestCmd = &cobra.Command{
	Use:   "digest",
	Short: "It will create or check digest file based on md5sum.",
	Long: `For creating a md5sum digest file, you must run this command at root of project or can 
	provide a root folder after the command. To verified is the same logic as create.`,
	Run: func(cmd *cobra.Command, args []string) {
		matches, err := zglob.Glob(args[0])

		if err != nil {
			color.Red("An error happend when tried to Glob this pattern %s, error: %s", args[0], err)
			return
		}

		if len(matches) <= 0 {
			color.Green("0 matches for pattern %s", args[0])
			return
		}
		color.Green("It was found matches: %b", len(matches))

		color.Cyan("Processing files...")
		done := make(chan bool)
		go digests(done, matches)
		<-done
		color.Cyan("We are done process...")

		if create {
			createFile()
			os.Exit(0)
		} else {
			if checkFile() {
				os.Exit(0)
			}
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(digestCmd)

	// Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// digestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// local flags which will only run when this command
	// is called directly, e.g.:
	digestCmd.Flags().BoolVarP(&create, "create", "c", false, "If provided, it will create a digest file rather than check")
}

func createFile() {
	f, err := os.Create("./checklist.chk")
	defer f.Close()
	if err != nil {
		color.Red("Unable to create a file")
		return
	}

	w := bufio.NewWriter(f)
	for _, line := range list {
		_, err1 := w.WriteString(line.String() + "\n")
		if err1 != nil {
			color.Red("Unable to write to a file")
		}
	}

	color.Green("Ok, writed. Check file ./checklist.chk")
	w.Flush()
}

func checkFile() bool {
	var status bool = true
	file, err := os.Open("./checklist.chk")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		partDigest := strings.Split(text, " ")
		md5 := partDigest[0]
		theFile := partDigest[1]
		for _, item := range list {
			if theFile == item.file {
				if item.digest != md5 && item.file != theFile {
					color.Red("The file %s is not the same as the digest check file. (new:%s, old:%s)", item.file, item.digest, md5)
					status = false
				} else {
					fmt.Printf("File %s is in good state\n", item.file)
				}
			}
		}
	}
	return status
}

func md5sum(file string) []byte {
	var ret []byte
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		return ret
	}
	input := strings.NewReader(string(data[:]))
	hash := md5.New()
	if _, err := io.Copy(hash, input); err != nil {
		log.Fatal(err)
		return ret
	}
	return hash.Sum(nil)
}

func digests(chDone chan bool, matches []string) {

	// Subprocess
	c1 := worker(matches)
	c2 := worker(matches)
	c3 := worker(matches)
	c4 := worker(matches)

	// Get Results..
	for v := range merge(c1, c2, c3, c4) {
		list = append(list, v)
	}

	chDone <- true
}

func worker(files []string) <-chan fileDigest {
	out := make(chan fileDigest)
	go func() {
		for _, file := range files {
			sum := md5sum(file)
			ret := fileDigest{file, fmt.Sprintf("%x", sum)}
			out <- ret
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan fileDigest) <-chan fileDigest {
	out := make(chan fileDigest)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan fileDigest) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (i fileDigest) String() string {
	return fmt.Sprintf("%s %s", i.digest, i.file)
}
