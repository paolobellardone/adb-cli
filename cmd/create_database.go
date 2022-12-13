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
	json "github.com/nwidger/jsoncolor"

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
		utils.PrintVerbose("OCI SDK for Go version: " + common.Version())

		// Aggiungere chiamate per scaricare il wallet e per stampare tutti i link di accesso al database
		//func (DatabaseClient) GenerateAutonomousDatabaseWallet
		//

		var adbName, _ = cmd.Flags().GetString("name")
		if len(adbName) > 14 {
			// If the name of the Autonomoous Database is longer than 14 chars then truncate it
			adbName = adbName[:14]
		}
		var ocpus int
		var storage int
		var database_type, _ = cmd.Flags().GetString("type")
		var adbType, isFree = utils.DecodeADBTypeForCreate(database_type)
		var enableOCPUAutoscaling bool
		var enableStorageAutoscaling bool
		var license_model database.CreateAutonomousDatabaseBaseLicenseModelEnum //database.CreateAutonomousDatabaseBaseLicenseModelEnum
		if !isFree {
			_tmp, _ := cmd.Flags().GetString("license-model")
			license_model = utils.DecodeLicenseModelForCreate(_tmp)
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
		/*
			dbClient, err := database.NewDatabaseClientWithConfigurationProvider(common.NewRawConfigurationProvider(ociConfig.tenancy, ociConfig.user, ociConfig.region, ociConfig.fingerprint, ociConfig.key_file, &ociConfig.pass_phrase))
			if err != nil {
				utils.PrintError("Error: " + err.Error())
				return
			}
		*/
		createADBDetails := database.CreateAutonomousDatabaseDetails{
			CompartmentId:                  &ociConfig.compartment_id,
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
		// Stampare se verbose la configurazione del DB da creare oppure stamparla una volta creato?
		s, _ := json.MarshalIndent(createADBDetails, "", "\t")
		utils.Print(string(s))

		/*
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
		*/
	},
}

func init() {
	createCmd.AddCommand(create_databaseCmd)

	create_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to create (required)")
	create_databaseCmd.MarkFlagRequired("name")
	create_databaseCmd.Flags().StringP("type", "t", "atpfree", "the type of the Autonomous Database to create - allowed values: atpfree, ajdfree, apexfree, adwfree, atp, ajd, apex, adw")
	create_databaseCmd.Flags().IntP("ocpus", "o", 1, "the number of OCPUs to allocate for the Autonomous Database - not used for Free Tier")
	create_databaseCmd.Flags().IntP("storage", "s", 1, "the size of storage in TB to allocate for the Autonomous Database - not used for Free Tier")
	create_databaseCmd.Flags().Bool("enable-ocpu-autoscaling", false, "enable autoscaling for OCPUs (max 3x the number of allocated OCPUs)")
	create_databaseCmd.Flags().Bool("enable-storage-autoscaling", false, "enable autoscaling for storage (max 3x the size of reserved storage)")
	create_databaseCmd.Flags().StringP("license-model", "l", "full", "the licensing model to use - allowed values: full, byol - not used for Free Tier")
}
