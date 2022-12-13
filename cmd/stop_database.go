/*
Copyright © 2022 PaoloB <paolo.bellardone@gmail.com>

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

// stop_databaseCmd stops the autonomous database specified by --name flag
var stop_databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Stop an Autonomous Database identified by the --name flag",
	Long:  "Stop an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("stop database command called")
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

		if adbInstance.LifecycleState == database.AutonomousDatabaseLifecycleStateAvailable {
			utils.PrintVerbose("The Autonomous Database " + adbName + " can be stopped")

			dbSADResponse, err := dbClient.StopAutonomousDatabase(context.Background(), database.StopAutonomousDatabaseRequest{AutonomousDatabaseId: common.String(adbOCID)})
			if err != nil {
				utils.PrintError("Error: " + err.Error())
				return
			}

			utils.PrintInfo("Stopping autonomous database " + adbName + ", please wait...")

			databaseStopping := true
			databaseStatus := dbSADResponse.AutonomousDatabase.LifecycleState

			for databaseStopping {
				adbInstance, err := utils.GetAutonomousDatabase(dbClient, adbOCID)
				if err != nil {
					utils.PrintError("Error: " + err.Error())
					return
				}

				if adbInstance.LifecycleState == database.AutonomousDatabaseLifecycleStateStopped {
					databaseStopping = false
					databaseStatus = adbInstance.LifecycleState
				} else {
					time.Sleep(15 * time.Second)
				}
			}

			utils.PrintInfo("The Autonomous Database " + *dbSADResponse.AutonomousDatabase.DbName + " was stopped, the current status is: " + string(databaseStatus))
		} else {
			utils.PrintWarning("The Autonomous Database " + adbName + " cannot be stopped because it is already stopped")
		}
	},
}

func init() {
	stopCmd.AddCommand(stop_databaseCmd)

	stop_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to stop (required)")
	stop_databaseCmd.MarkFlagRequired("name")
}
