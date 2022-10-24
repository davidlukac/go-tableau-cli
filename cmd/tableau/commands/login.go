package commands

/*
Copyright © 2022 David Lukac <1215290+davidlukac@users.noreply.github.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

import (
	"fmt"
	"github.com/davidlukac/go-tableau-cli/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Performs login and prints site ID and user ID used for the login.",
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("tableau_url")
		username := viper.GetString("tableau_username")
		password := viper.GetString("tableau_password")
		log.Infof("Logging in into %s as %s...", url, username)
		token, siteId, err := internal.Login(url, username, password)
		if err != nil {
			log.Errorf("Command failed: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully logged into %s (site ID %s); token is %s\n", url, siteId, token)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
