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

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"time"
)

var get bool
var url string

// webdiffCmd represents the webdiff command
var webdiffCmd = &cobra.Command{
	Use:   "webdiff",
	Short: "Security related to change the visual of a webpage",
	Long:  `It use chrome headless and image diff to compare which part of a webpage was different...`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		var buf []byte
		color.Green(url)
		err := chromedp.Run(ctx, screenshot(url, &buf))
		if err != nil {
			log.Fatal(err)
		}

		// save the screenshot to disk
		if err = ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webdiffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webdiffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	webdiffCmd.Flags().BoolVarP(&get, "get", "g", false, "Get a new Image from Url.")
	webdiffCmd.Flags().StringVarP(&url, "url", "u", "", "Get screenshot from this url.")
	webdiffCmd.MarkFlagRequired("url")
}

func screenshot(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(time.Duration(3) * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}
