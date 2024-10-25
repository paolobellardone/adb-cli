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
	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// delete_configCmd deletes the configuration file specified by --config global flag
var delete_configCmd = &cobra.Command{
	Use:   "config",
	Short: "Delete the configuration file specified by --config global flag",
	Long:  "Delete the configuration file specified by --config global flag",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("delete config command called")
		utils.PrintVerbose("")

		// TODO: implement the feature
		// TODO: delete the config file or the config template?
		utils.PrintError("Feature not impleleted yet!")
	},
}

func init() {
	deleteCmd.AddCommand(delete_configCmd)
}
