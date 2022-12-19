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

	json "github.com/nwidger/jsoncolor"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// update_databaseCmd updates the parameters of the autonomous database specified by --name flag
var update_databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Update the parameters of an Autonomous Database identified by the --name flag",
	Long:  "Update the parameters of an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("update database command called")
		utils.PrintVerbose("OCI SDK for Go version: " + common.Version())

		var adbName, _ = cmd.Flags().GetString("name")
		if len(adbName) > 14 {
			// If the name of the Autonomoous Database is longer than 14 chars then truncate it
			adbName = adbName[:14]
		}
		var ocpus int
		var storage int
		var database_type, _ = cmd.Flags().GetString("type")
		var adbType, isFree = utils.DecodeADBTypeForUpdate(database_type)
		var enableOCPUAutoscaling bool
		var enableStorageAutoscaling bool
		var license_model database.UpdateAutonomousDatabaseDetailsLicenseModelEnum
		if !isFree {
			_tmp, _ := cmd.Flags().GetString("license-model")
			license_model = utils.DecodeLicenseModelForUpdate(_tmp)
			ocpus, _ = cmd.Flags().GetInt("ocpus")
			storage, _ = cmd.Flags().GetInt("storage")
			enableOCPUAutoscaling, _ = cmd.Flags().GetBool("enable-ocpu-autoscaling")
			enableStorageAutoscaling, _ = cmd.Flags().GetBool("enable-storage-autoscaling")
		} else {
			// If the Autonomous Database is a Free Tier, the OCPUs are limited to 1 and the Storage is limited to 20 GB
			// TODO: is it really needed???
			ocpus = 1
			storage = 1
		}

		dbClient, err := database.NewDatabaseClientWithConfigurationProvider(common.NewRawConfigurationProvider(ociConfig.tenancy, ociConfig.user, ociConfig.region, ociConfig.fingerprint, ociConfig.key_file, &ociConfig.pass_phrase))
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		updateADBDetails := database.UpdateAutonomousDatabaseDetails{
			DbName:                         common.String(adbName),
			CpuCoreCount:                   common.Int(ocpus),
			DataStorageSizeInTBs:           common.Int(storage),
			IsFreeTier:                     &isFree,
			AdminPassword:                  &ociConfig.database_password,
			DisplayName:                    common.String(adbName),
			IsAutoScalingEnabled:           common.Bool(enableOCPUAutoscaling),
			IsAutoScalingForStorageEnabled: common.Bool(enableStorageAutoscaling),
			DbWorkload:                     adbType,
			LicenseModel:                   license_model,
			// TODO: Add other parameters?
		}

		// TODO: print the configured options before and/or after the update of the ADB???
		s, _ := json.MarshalIndent(updateADBDetails, "", "\t")
		utils.Print(string(s))

		updateADBResponse, err := dbClient.UpdateAutonomousDatabase(
			context.Background(),
			database.UpdateAutonomousDatabaseRequest{
				UpdateAutonomousDatabaseDetails: updateADBDetails,
			})
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		utils.PrintInfo("The Autonomous Database " + *updateADBResponse.AutonomousDatabase.DbName + " status is currently: " + string(updateADBResponse.AutonomousDatabase.LifecycleState))
		if updateADBResponse.AutonomousDatabase.LifecycleState == database.AutonomousDatabaseLifecycleStateUpdating {
			utils.PrintVerbose("Updating autonomous database " + adbName + ", please wait...")

			databaseUpdating := true
			databaseStatus := updateADBResponse.AutonomousDatabase.LifecycleState

			for databaseUpdating {
				adbInstance, err := utils.GetAutonomousDatabase(dbClient, *updateADBResponse.AutonomousDatabase.Id)
				if err != nil {
					utils.PrintError("Error: " + err.Error())
					return
				}

				if adbInstance.LifecycleState == database.AutonomousDatabaseLifecycleStateAvailable {
					databaseUpdating = false
					databaseStatus = adbInstance.LifecycleState
				} else {
					time.Sleep(15 * time.Second)
				}
			}

			utils.PrintInfo("The Autonomous Database " + *updateADBResponse.AutonomousDatabase.DbName + " was updated. The current status is: " + string(databaseStatus))
		} else {
			utils.PrintError("Error during updating of the Autonomous Database " + adbName + ". Please check the OCI console for errors...")
			return
		}
	},
}

func init() {
	updateCmd.AddCommand(update_databaseCmd)

	update_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to update (required)")
	update_databaseCmd.MarkFlagRequired("name")
	update_databaseCmd.Flags().StringP("type", "t", "atpfree", "the type of the Autonomous Database to update - allowed values: atpfree, ajdfree, apexfree, adwfree, atp, ajd, apex, adw")
	update_databaseCmd.Flags().IntP("ocpus", "o", 1, "the number of OCPUs to allocate for the Autonomous Database - not used for Free Tier")
	update_databaseCmd.Flags().IntP("storage", "s", 1, "the size of storage in TB to allocate for the Autonomous Database - not used for Free Tier")
	update_databaseCmd.Flags().Bool("enable-ocpu-autoscaling", false, "enable autoscaling for OCPUs (max 3x the number of allocated OCPUs)")
	update_databaseCmd.Flags().Bool("enable-storage-autoscaling", false, "enable autoscaling for storage (max 3x the size of reserved storage)")
	update_databaseCmd.Flags().StringP("license-model", "l", "full", "the licensing model to use - allowed values: full, byol - not used for Free Tier")
}
