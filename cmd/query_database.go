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
	"strconv"

	json "github.com/nwidger/jsoncolor"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// query_databaseCmd queries the autonomous database specified by --name flag
var query_databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Query an Autonomous Database identified by the --name flag",
	Long:  "Query an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("query database command called")
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

		if full, _ := cmd.Flags().GetBool("full"); full {
			s, _ := json.MarshalIndent(adbInstance, "", "\t")
			utils.Print(string(s))
		} else {
			utils.PrintKV("Database status             : ", string(adbInstance.LifecycleState))
			utils.PrintKV("Database name               : ", *adbInstance.DbName)
			utils.PrintKV("Display name                : ", *adbInstance.DisplayName)
			utils.PrintKV("Database version            : ", *adbInstance.DbVersion)
			utils.PrintKV("Preview version             : ", strconv.FormatBool(*adbInstance.IsPreview))
			utils.PrintKV("OCPUs                       : ", strconv.Itoa(*adbInstance.CpuCoreCount))
			utils.PrintKV("Storage (GB)                : ", strconv.Itoa(*adbInstance.DataStorageSizeInGBs))
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
	queryCmd.AddCommand(query_databaseCmd)

	query_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to inspect (required)")
	query_databaseCmd.MarkFlagRequired("name")
	query_databaseCmd.Flags().BoolP("full", "f", false, "print all the attributes of the Autonomous Database")
}
