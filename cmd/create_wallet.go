/*
Copyright © 2022, 2024s PaoloB

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
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// This is the default password that will be used for a wallet if not specified by the --password flag
const default_wallet_password string = "WElcome1234__"

// create_walletCmd creates the wallet for the autonomous database specified by --name flag
var create_walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Create the wallet for an Autonomous Database identified by the --name flag",
	Long:  "Create the wallet for an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("create wallet command called")
		utils.PrintVerbose("")

		var adbName, _ = cmd.Flags().GetString("name")
		var walletPassword, _ = cmd.Flags().GetString("password")
		if walletPassword == "" {
			utils.PrintWarning("--password flag not specified, using default password")
			walletPassword = default_wallet_password
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

		walletName, err := utils.CreateAutonomousDatabaseWallet(dbClient, adbName, adbOCID, walletPassword)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		} else {
			utils.PrintInfo("The wallet " + walletName + " for Autonomous Database " + adbName + " was successfully created.")
		}
	},
}

func init() {
	createCmd.AddCommand(create_walletCmd)

	create_walletCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database for which create the wallet (required)")
	create_walletCmd.MarkFlagRequired("name")
	create_walletCmd.Flags().String("password", "", "the password that protects the wallet")
}
