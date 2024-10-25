/*
Copyright © 2022, 2024 PaoloB

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

	//json "github.com/nwidger/jsoncolor"

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
		utils.PrintVerbose("")

		var adbName string
		var databaseType string
		var ocpus int
		var storage int
		var enableOCPUAutoscaling bool
		var enableStorageAutoscaling bool
		//var disableOCPUAutoscaling bool
		//var disableStorageAutoscaling bool
		var licenseModel string

		// Semplificare cosa si può modificare, sicuramente non il tipo (adw/atp ecc) mettere solo le configurazioni di base cpu/storage ecc
		// E' possibile solo fare una operazione per volta!!! modificare
		// per cambiare la licenza a byol è necessario specificare il tipo di licenza, anche per create!!! (Standard o Enterprise)
		var updateADBDetails database.UpdateAutonomousDatabaseDetails
		var adbType interface{}
		var isFree bool

		adbName, _ = cmd.Flags().GetString("name")
		if len(adbName) > 14 {
			// If the name of the Autonomous Database is longer than 14 chars then truncate it
			adbName = adbName[:14]
		}

		if cmd.Flags().Changed("type") {
			databaseType, _ = cmd.Flags().GetString("type")
			adbType, isFree = utils.DecodeADBType(databaseType, "update")
			updateADBDetails.DbWorkload = adbType.(database.UpdateAutonomousDatabaseDetailsDbWorkloadEnum)
			updateADBDetails.IsFreeTier = common.Bool(isFree)
		}
		if cmd.Flags().Changed("ocpus") {
			ocpus, _ = cmd.Flags().GetInt("ocpus")
			updateADBDetails.CpuCoreCount = common.Int(ocpus)
		}
		if cmd.Flags().Changed("storage") {
			storage, _ = cmd.Flags().GetInt("storage")
			updateADBDetails.DataStorageSizeInTBs = common.Int(storage)
		}
		if cmd.Flags().Changed("enable-ocpu-autoscaling") {
			enableOCPUAutoscaling, _ = cmd.Flags().GetBool("enable-ocpu-autoscaling")
			updateADBDetails.IsAutoScalingEnabled = common.Bool(enableOCPUAutoscaling)
		}
		if cmd.Flags().Changed("enable-storage-autoscaling") {
			enableStorageAutoscaling, _ = cmd.Flags().GetBool("enable-storage-autoscaling")
			updateADBDetails.IsAutoScalingEnabled = common.Bool(enableStorageAutoscaling)
		}
		if cmd.Flags().Changed("disable-ocpu-autoscaling") {
			updateADBDetails.IsAutoScalingEnabled = common.Bool(false)
		}
		if cmd.Flags().Changed("disable-storage-autoscaling") {
			updateADBDetails.IsAutoScalingEnabled = common.Bool(false)
		}
		if cmd.Flags().Changed("license-model") {
			licenseModel, _ = cmd.Flags().GetString("license-model")
			var _t1, _t2 = utils.DecodeLicenseModel(licenseModel, "update")
			updateADBDetails.LicenseModel = _t1.(database.UpdateAutonomousDatabaseDetailsLicenseModelEnum)
			updateADBDetails.DatabaseEdition = _t2.(database.AutonomousDatabaseSummaryDatabaseEditionEnum)
		}

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

		// TODO: print the configured options before and/or after the update of the ADB???
		//s, _ := json.MarshalIndent(updateADBDetails, "", "\t")
		//utils.Print(string(s))

		updateADBResponse, err := dbClient.UpdateAutonomousDatabase(
			context.Background(),
			database.UpdateAutonomousDatabaseRequest{
				AutonomousDatabaseId:            common.String(adbOCID),
				UpdateAutonomousDatabaseDetails: updateADBDetails,
			})
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		utils.PrintInfo("The Autonomous Database " + *updateADBResponse.AutonomousDatabase.DbName + " status is currently: " + string(updateADBResponse.AutonomousDatabase.LifecycleState))
		if updateADBResponse.AutonomousDatabase.LifecycleState == database.AutonomousDatabaseLifecycleStateUpdating ||
			updateADBResponse.AutonomousDatabase.LifecycleState == database.AutonomousDatabaseLifecycleStateScaleInProgress {
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
	update_databaseCmd.Flags().Bool("disable-ocpu-autoscaling", false, "disable autoscaling for OCPUs")
	update_databaseCmd.Flags().Bool("disable-storage-autoscaling", false, "disable autoscaling for storage")
	update_databaseCmd.Flags().StringP("license-model", "l", "full", "the licensing model to use - allowed values: full, byolee, byolse - not used for Free Tier")
}
