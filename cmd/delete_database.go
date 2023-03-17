/*
Copyright © 2022 PaoloB

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
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// delete_databaseCmd deletes the autonomous database specified by --name flag
var delete_databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Delete an Autonomous Database identified by the --name flag",
	Long:  "Delete an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("delete database command called")
		utils.PrintVerbose("")

		var adbName, _ = cmd.Flags().GetString("name")

		dbClient, err := database.NewDatabaseClientWithConfigurationProvider(common.NewRawConfigurationProvider(ociConfig.tenancy, ociConfig.user, ociConfig.region, ociConfig.fingerprint, ociConfig.key_file, &ociConfig.pass_phrase))
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		adbOCID, err := utils.GetAutonomousDatabaseOCID(dbClient, adbName, ociConfig.compartment_id)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		// The delete command shoiuld always be confirmed!
		fmt.Print("Please confirm by entering the name of the Autonomous Database you want to delete (" + adbName + "): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if scanner.Text() == adbName {
			color.Set(color.FgHiYellow)
			fmt.Print("You asked to delete the Autonomous Database named " + scanner.Text() + ", confirm? (Y/n) ")
			color.Unset()
			scanner.Scan()
			if scanner.Text() != "Y" {
				utils.PrintError("The termination of the Autonomous Database " + adbName + " was ABORTED by the user!")
			} else {
				// Send the request using the service client
				_, err := dbClient.DeleteAutonomousDatabase(context.Background(), database.DeleteAutonomousDatabaseRequest{AutonomousDatabaseId: common.String(adbOCID)})
				if err != nil {
					utils.PrintError("Error: " + err.Error())
					return
				}

				utils.PrintInfo("Terminating autonomous database " + adbName + ", please wait...")

				adbInstance, err := utils.GetAutonomousDatabase(dbClient, adbOCID)
				if err != nil {
					utils.PrintError("Error: " + err.Error())
					return
				}

				databaseTerminating := true
				databaseStatus := adbInstance.LifecycleState

				for databaseTerminating {
					adbInstance, err := utils.GetAutonomousDatabase(dbClient, adbOCID)
					if err != nil {
						utils.PrintError("Error: " + err.Error())
						return
					}

					if adbInstance.LifecycleState == database.AutonomousDatabaseLifecycleStateTerminated {
						databaseTerminating = false
						databaseStatus = adbInstance.LifecycleState
					} else {
						time.Sleep(15 * time.Second)
					}
				}
				utils.PrintInfo("The Autonomous Database " + *adbInstance.DbName + " was terminated, the current status is: " + string(databaseStatus))
			}
		} else {
			utils.PrintError("Error: The specified Autonomous Database name does not match")
			return
		}

		// Delete the Autonomous Database wallet for the specified database (it should be there if it was created by adb-cli)
		err = utils.DeleteAutonmousDatabaseWallet(adbName)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		} else {
			utils.PrintInfo("The wallet for Autonomous Database " + adbName + " was successfully deleted.")
		}
	},
}

func init() {
	deleteCmd.AddCommand(delete_databaseCmd)

	delete_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to delete (required)")
	delete_databaseCmd.MarkFlagRequired("name")
}
