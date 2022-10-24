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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

const (
	SiteRoleFlagName = "site-role"
	FromYamlFlagName = "from-yaml"
)

var (
	siteRoleFlag string
	fromYamlFlag string
)

// updateUserCmd represents the updateUser command
var updateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Update existing user properties",
	Long: fmt.Sprintf(`
Update can be done with provided username argument and %s flag, or
with list of users and roles in YAML file with usernames and roles formatted as 
- username: john.smith
  role: Explorer
`, SiteRoleFlagName),
	PreRun: internal.LoggingSetup,
	Args:   cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && len(siteRoleFlag) > 0 {
			// Updating one user with given role
			username := args[0]

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

			user, err := t.UpdateUserSiteRole(username, siteRoleFlag)
			if err != nil {
				log.Errorf("Command failed: %s", err)
				os.Exit(1)
			}

			fmt.Printf("%s (%s) - %s\n", user.Username, user.ID, user.Role)
		} else if len(args) == 0 && len(fromYamlFlag) > 0 {
			fileInfo, err := os.Stat(fromYamlFlag)
			if os.IsNotExist(err) {
				log.Errorf("Provided path to YAML file is not valid: %s", err)
				os.Exit(1)
			} else if fileInfo.IsDir() {
				log.Errorf("Provided path to YAML file is not valid: it's a directory")
				os.Exit(1)
			}
			log.Infof("Updating user roles from file %s", fromYamlFlag)

			yamlFile, err := ioutil.ReadFile(fromYamlFlag)
			if err != nil {
				log.Printf("Couldn't read YAML file: %v ", err)
			}

			var users []*internal.User

			err = yaml.Unmarshal(yamlFile, &users)
			if err != nil {
				log.Fatalf("Unmarshal: %v", err)
			}

			log.Debugf("Loaded %d user from file", len(users))

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

			alreadySame := 0
			updated := 0
			notExists := 0
			errored := 0

			for idx, user := range users {
				actual, err := t.GetUser(user.Username)
				if err != nil {
					errored++
					log.Errorf("[%d/%d] Failed to fetch user %s - skipping update", idx, len(users), user.Username)
				} else {
					if actual.Exists {
						if strings.ToLower(actual.Role) == strings.ToLower(user.Role) {
							alreadySame++
							log.Infof("[%d/%d] User %s already has role %s", idx, len(users), user.Username, actual.Role)
						} else {
							log.Infof("[%d/%d] Updating user %s from role %s to role %s...", idx, len(users), user.Username, actual.Role, user.Role)
							updateUser, err := t.UpdateUserSiteRole(user.Username, user.Role)
							if err != nil {
								errored++
								log.Errorf("Failed to update user %s to role %s: %v", user.Username, user.Role, err)
							} else {
								updated++
								log.Infof("... user %s update to role %s", user.Username, updateUser.Role)
							}
						}
					} else {
						notExists++
						log.Warnf("[%d/%d] User %s doesn't not existing! Skipping...", idx, len(users), user.Username)
					}
				}
			}
			fmt.Printf("\nAlready same role: %d\nUpdated: %d\nNot found: %d\nError: %d\n", alreadySame, updated, notExists, errored)

		} else {
			_ = cmd.Help()
		}
	},
}

func init() {
	updateCmd.AddCommand(updateUserCmd)

	updateUserCmd.Flags().StringVar(&siteRoleFlag, SiteRoleFlagName, "", "User's role")
	updateUserCmd.Flags().StringVar(&fromYamlFlag, FromYamlFlagName, "", "Path to YAML file with username(s) and role(s)")
}
