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
	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// createCmd represents the create command, the allowed resources for this command are: config, database, wallet
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an OCI resource and/or supporting configuration - allowed resources: config, database, wallet",
	Long:  "Create an OCI resource and/or supporting configuration.\nThe allowed resources are: config, database, wallet",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintError("Please specify a resource. Allowed values: config, database, wallet")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
