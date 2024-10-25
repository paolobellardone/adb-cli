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
	"strings"
	"time"

	//json "github.com/nwidger/jsoncolor"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// create_databaseCmd creates the autonomous database specified by --name flag
var create_databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Create an Autonomous Database identified by the --name flag",
	Long:  "Create an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("create database command called")
		utils.PrintVerbose("")

		// TODO: print the connection URL for the database and/or the connection strings???

		var adbName string
		var databaseType string
		var computeModel string
		var computeUnits int
		var storage int
		var enableOCPUAutoscaling bool
		var enableStorageAutoscaling bool
		var licenseModel string
		var createADBDetails database.CreateAutonomousDatabaseDetails
		var adbType interface{}
		var isFree bool

		// Evaluate first the parameters specified inside the configuration file and then the command line
		createADBDetails.CompartmentId = common.String(ociConfig.compartment_id)
		createADBDetails.AdminPassword = common.String(ociConfig.database_password)

		adbName, _ = cmd.Flags().GetString("name")
		if len(adbName) > 14 {
			// If the name of the Autonomous Database is longer than 14 chars then truncate it
			adbName = adbName[:14]
		}
		createADBDetails.DbName = common.String(adbName)
		createADBDetails.DisplayName = common.String(adbName)

		databaseType, _ = cmd.Flags().GetString("type")
		adbType, isFree = utils.DecodeADBType(databaseType, "create")
		createADBDetails.DbWorkload = adbType.(database.CreateAutonomousDatabaseBaseDbWorkloadEnum)
		createADBDetails.IsFreeTier = common.Bool(isFree)

		computeModel, _ = cmd.Flags().GetString("compute-model")
		if !isFree {
			createADBDetails.ComputeModel = database.CreateAutonomousDatabaseBaseComputeModelEnum(strings.ToUpper(computeModel))
		} else {
			createADBDetails.ComputeModel = database.CreateAutonomousDatabaseBaseComputeModelEnum("ECPU")
		}

		computeUnits, _ = cmd.Flags().GetInt("compute-units")
		if !isFree {
			if strings.ToUpper(computeModel) == "ECPU" && (computeUnits%2) != 0 {
				utils.PrintError("Error: when using ECPU compute model the number of compute units should be a multiple of 2")
				return
			} else {
				createADBDetails.ComputeCount = common.Float32(float32(computeUnits))
			}
		} else {
			createADBDetails.ComputeCount = common.Float32(1)
		}

		storage, _ = cmd.Flags().GetInt("storage")
		if !isFree {
			createADBDetails.DataStorageSizeInTBs = common.Int(storage)
		} else {
			createADBDetails.DataStorageSizeInTBs = common.Int(1)
		}

		enableOCPUAutoscaling, _ = cmd.Flags().GetBool("enable-ocpu-autoscaling")
		if !isFree {
			createADBDetails.IsAutoScalingEnabled = common.Bool(enableOCPUAutoscaling)
		}
		enableStorageAutoscaling, _ = cmd.Flags().GetBool("enable-storage-autoscaling")
		if !isFree {
			createADBDetails.IsAutoScalingForStorageEnabled = common.Bool(enableStorageAutoscaling)
		}

		licenseModel, _ = cmd.Flags().GetString("license-model")
		if !isFree {
			_t1, _t2 := utils.DecodeLicenseModel(licenseModel, "create")
			createADBDetails.LicenseModel = _t1.(database.CreateAutonomousDatabaseBaseLicenseModelEnum)
			// The DatabaseEdition parameter is needed only if the LicenseModel is BYOL
			if _t2 != nil {
				createADBDetails.DatabaseEdition = _t2.(database.AutonomousDatabaseSummaryDatabaseEditionEnum)
			}
		}

		if ociConfig.database_password == "" {
			utils.PrintError("Error: database_password paramenter in config file can not be empty")
			return
		}

		dbClient, err := database.NewDatabaseClientWithConfigurationProvider(common.NewRawConfigurationProvider(ociConfig.tenancy, ociConfig.user, ociConfig.region, ociConfig.fingerprint, ociConfig.key_file, &ociConfig.pass_phrase))
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		// TODO: print the configured options before and/or after the creation of the ADB???
		//s, _ := json.MarshalIndent(createADBDetails, "", "\t")
		//utils.Print(string(s))

		createADBResponse, err := dbClient.CreateAutonomousDatabase(
			context.Background(),
			database.CreateAutonomousDatabaseRequest{
				CreateAutonomousDatabaseDetails: createADBDetails,
			})
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		utils.PrintInfo("The Autonomous Database " + *createADBResponse.AutonomousDatabase.DbName + " status is currently: " + string(createADBResponse.AutonomousDatabase.LifecycleState))
		if createADBResponse.AutonomousDatabase.LifecycleState == database.AutonomousDatabaseLifecycleStateProvisioning {
			utils.PrintVerbose("Provisioning autonomous database " + adbName + ", please wait...")

			databaseProvisioning := true
			databaseStatus := createADBResponse.AutonomousDatabase.LifecycleState

			for databaseProvisioning {
				adbInstance, err := utils.GetAutonomousDatabase(dbClient, *createADBResponse.AutonomousDatabase.Id)
				if err != nil {
					utils.PrintError("Error: " + err.Error())
					return
				}

				if adbInstance.LifecycleState == database.AutonomousDatabaseLifecycleStateAvailable {
					databaseProvisioning = false
					databaseStatus = adbInstance.LifecycleState
				} else {
					time.Sleep(15 * time.Second)
				}
			}

			utils.PrintInfo("The Autonomous Database " + *createADBResponse.AutonomousDatabase.DbName + " was provisioned. The current status is: " + string(databaseStatus))
		} else {
			utils.PrintError("Error during creation of the Autonomous Database " + adbName + ". Please check the OCI console for errors...")
			return
		}

		// The default password for wallet is the same of the ADMIN user
		walletName, err := utils.CreateAutonomousDatabaseWallet(dbClient, adbName, *createADBResponse.AutonomousDatabase.Id, ociConfig.database_password)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		} else {
			utils.PrintInfo("The wallet " + walletName + " for Autonomous Database " + adbName + " was successfully created.")
			utils.PrintInfo("The password for the wallet is the same of ADMIN user.")
		}
	},
}

func init() {
	createCmd.AddCommand(create_databaseCmd)

	create_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to create (required)")
	create_databaseCmd.MarkFlagRequired("name")
	create_databaseCmd.Flags().StringP("type", "t", "atpfree", "the type of the Autonomous Database to create - allowed values: atpfree, ajdfree, apexfree, adwfree, atp, ajd, apex, adw")
	create_databaseCmd.Flags().StringP("compute-model", "m", "ECPU", "the compute model to be used to create the Autonomous Database - allowed values: ECPU, OCPU -- not used for Free Tier")
	create_databaseCmd.Flags().IntP("compute-units", "u", 2, "the number of compute units to allocate for the Autonomous Database -- not used for Free Tier")
	create_databaseCmd.Flags().IntP("storage", "s", 1, "the size of storage in TB to allocate for the Autonomous Database -- not used for Free Tier")
	create_databaseCmd.Flags().Bool("enable-storage-autoscaling", false, "enable autoscaling for storage (max 3x the size of reserved storage) -- not used for Free Tier")
	create_databaseCmd.Flags().StringP("license-model", "l", "full", "the licensing model to use - allowed values: full, byolee, byolse -- not used for Free Tier")
}
