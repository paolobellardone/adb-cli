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
	"strconv"

	json "github.com/nwidger/jsoncolor"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// inspect_databaseCmd inspects the autonomous database specified by --name flag
var inspect_databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Inspect an Autonomous Database identified by the --name flag",
	Long:  "Inspect an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("inspect database command called")
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

		adbInstance, err := utils.GetAutonomousDatabase(dbClient, adbOCID)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		}

		if full, _ := cmd.Flags().GetBool("full"); full {
			s, _ := json.MarshalIndent(adbInstance, "", "\t")
			utils.Print(string(s))
		} else {
			utils.PrintKV("Database status             : ", string(adbInstance.LifecycleState))
			utils.PrintKV("Database name               : ", *adbInstance.DbName)
			utils.PrintKV("Display name                : ", *adbInstance.DisplayName)
			utils.PrintKV("Database version            : ", *adbInstance.DbVersion)
			utils.PrintKV("Preview version             : ", strconv.FormatBool(*adbInstance.IsPreview))
			utils.PrintKV(string(adbInstance.ComputeModel)+"s                       : ", strconv.FormatFloat(float64(*adbInstance.ComputeCount), 'f', -1, 32))
			utils.PrintKV("Configured storage (GB)     : ", strconv.Itoa(*adbInstance.DataStorageSizeInGBs))
			utils.PrintKV("Used storage (GB)           : ", strconv.FormatFloat(float64(*adbInstance.ActualUsedDataStorageSizeInTBs*1024), 'f', -1, 32))
			utils.PrintKV("Workload type               : ", string(adbInstance.DbWorkload))
			if *adbInstance.IsFreeTier {
				utils.PrintKV("Free tier                   : ", strconv.FormatBool(*adbInstance.IsFreeTier))
			} else {
				utils.PrintKV("License model               : ", string(adbInstance.LicenseModel))
			}
			utils.PrintKV("OCPU auto-scaling enabled   : ", strconv.FormatBool(*adbInstance.IsAutoScalingEnabled))
			utils.PrintKV("Storage auto-scaling enabled: ", strconv.FormatBool(*adbInstance.IsAutoScalingForStorageEnabled))
			utils.Print("--------------")
			// TODO: print only the details needed to create or update an ADB, add other details??? (Data Guard???)
			utils.PrintInfo("Connections URLs")
			utils.PrintKV("Database Actions            : ", *adbInstance.ConnectionUrls.SqlDevWebUrl)
			// The following two are not needed anymore because all is available in Database Actions
			//utils.PrintKV("Oracle Application Express  : ", *adbInstance.ConnectionUrls.ApexUrl)
			//utils.PrintKV("Oracle Graph Studio         : ", *adbInstance.ConnectionUrls.GraphStudioUrl)
		}
	},
}

func init() {
	inspectCmd.AddCommand(inspect_databaseCmd)

	inspect_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to inspect (required)")
	inspect_databaseCmd.MarkFlagRequired("name")
	inspect_databaseCmd.Flags().BoolP("full", "f", false, "print all the attributes of the Autonomous Database")
}
