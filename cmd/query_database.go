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
	"fmt"
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
			fmt.Print("Database name               : ")
			utils.PrintInfo(*adbInstance.DbName)
			fmt.Print("Display name                : ")
			utils.PrintInfo(*adbInstance.DisplayName)
			fmt.Print("Databasse version           : ")
			utils.PrintInfo(*adbInstance.DbVersion)
			fmt.Print("Preview version             : ")
			utils.PrintInfo(strconv.FormatBool(*adbInstance.IsPreview))
			fmt.Print("OCPUs                       : ")
			utils.PrintInfo(strconv.Itoa(*adbInstance.CpuCoreCount))
			fmt.Print("Storage (GB)                : ")
			utils.PrintInfo(strconv.Itoa(*adbInstance.DataStorageSizeInGBs))
			fmt.Print("Workload type               : ")
			utils.PrintInfo(string(adbInstance.DbWorkload))
			if *adbInstance.IsFreeTier {
				fmt.Print("Free tier                   : ")
				utils.PrintInfo(strconv.FormatBool(*adbInstance.IsFreeTier))
			} else {
				fmt.Print("License model               : ")
				utils.PrintInfo(string(adbInstance.LicenseModel))
			}
			fmt.Print("OCPU auto-scaling enabled   : ")
			utils.PrintInfo(strconv.FormatBool(*adbInstance.IsAutoScalingEnabled))
			fmt.Print("Storage auto-scaling enabled: ")
			utils.PrintInfo(strconv.FormatBool(*adbInstance.IsAutoScalingForStorageEnabled))
			// TODO: print only the details needed to create or update an ADB, add other details??? (Data Guard???)
		}
	},
}

func init() {
	queryCmd.AddCommand(query_databaseCmd)

	query_databaseCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database to inspect (required)")
	query_databaseCmd.MarkFlagRequired("name")
	// TODO: add the flags to print all the info or only the essential - use global verbosity flag???
	query_databaseCmd.Flags().BoolP("full", "f", false, "print all the attributes of the Autonomous Database")
}
