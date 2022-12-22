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
	"github.com/oracle/oci-go-sdk/v65/common"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// delete_walletCmd deletes the wallet for the autonomous database specified by --name flag
var delete_walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Delete the wallet for an Autonomous Database identified by the --name flag",
	Long:  "Delete the wallet for an Autonomous Database identified by the --name flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("delete wallet command called")
		utils.PrintVerbose("OCI SDK for Go version: " + common.Version())

		var adbName, _ = cmd.Flags().GetString("name")

		err := utils.DeleteAutonmousDatabaseWallet(adbName)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		} else {
			utils.PrintInfo("The wallet for Autonomous Database " + adbName + " was successfully deleted.")
		}
	},
}

func init() {
	deleteCmd.AddCommand(delete_walletCmd)

	delete_walletCmd.Flags().StringP("name", "n", "", "the name of the Autonomous Database associated with the wallet to delete (required)")
	delete_walletCmd.MarkFlagRequired("name")
}
