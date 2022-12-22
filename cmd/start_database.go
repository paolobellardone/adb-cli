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
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// start_databaseCmd starts the autonomous database specified by --name flag
var start_databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Start an Autonomous Database identified by the --name flag",
	Long:  "Start an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("start database command called")
		utils.PrintVerbose("OCI SDK for Go version: " + common.Version())

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

		adbInstance, err := utils.GetAutonomousDatabase(dbClient, adbOCID)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		utils.PrintInfo("The current status of the Autonomous Database " + *adbInstance.DbName + " is: " + string(adbInstance.LifecycleState))

		if adbInstance.LifecycleState == database.AutonomousDatabaseLifecycleStateStopped {
			utils.PrintVerbose("The Autonomous Database " + adbName + " can be started")

			dbSADResponse, err := dbClient.StartAutonomousDatabase(context.Background(), database.StartAutonomousDatabaseRequest{AutonomousDatabaseId: common.String(adbOCID)})
			if err != nil {
				utils.PrintError("Error: " + err.Error())
				return
			}

			utils.PrintInfo("Starting autonomous database " + adbName + ", please wait...")

			databaseStarting := true
			databaseStatus := dbSADResponse.AutonomousDatabase.LifecycleState

			for databaseStarting {
				adbInstance, err := utils.GetAutonomousDatabase(dbClient, adbOCID)
				if err != nil {
					utils.PrintError("Error: " + err.Error())
					return
				}

				if adbInstance.LifecycleState == database.AutonomousDatabaseLifecycleStateAvailable {
					databaseStarting = false
					databaseStatus = adbInstance.LifecycleState
				} else {
					time.Sleep(15 * time.Second)
				}
			}

			utils.PrintInfo("The Autonomous Database " + *dbSADResponse.AutonomousDatabase.DbName + " was started, the current status is: " + string(databaseStatus))
		} else {
			utils.PrintWarning("The Autonomous Database " + adbName + " cannot be started because it is already running")
		}
	},
}

func init() {
	startCmd.AddCommand(start_databaseCmd)

	start_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to start (required)")
	start_databaseCmd.MarkFlagRequired("name")
}
