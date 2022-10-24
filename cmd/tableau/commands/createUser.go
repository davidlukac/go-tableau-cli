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

// createUserCmd represents the createUser command
var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		fmt.Printf("createUser called with %s\n", username)

		url := viper.GetString("tableau_url")
		apiUsername := viper.GetString("tableau_username")
		apiPassword := viper.GetString("tableau_password")
		token, siteId, err := internal.Login(url, apiUsername, apiPassword)
		if err != nil {
			log.Errorf("Failed to log in: %s", err)
			os.Exit(1)
		}

		t := internal.Tableau{
			BaseURL: url,
			Token:   token,
			SiteID:  siteId,
		}

		user, err := t.CreateUser(internal.User{
			Exists:   false,
			Username: username,
			Role:     internal.DefaultRole,
		})
		if err != nil {
			if user.Exists {
				fmt.Printf("User %s already exists!\n", user.Username)
			} else {
				log.Errorf("Command failed: %s", err)
			}
			os.Exit(1)
		}

		fmt.Printf("User %s created with ID %s (%v)\n", user.Username, user.ID, user)
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)
}
