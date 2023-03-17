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
	"errors"
	"os"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/oracle/oci-go-sdk/v65/common"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

const cliVersion string = "0.2.4"

type ociConfigT struct {
	user               string
	key_file           string
	pass_phrase        string
	fingerprint        string
	tenancy            string
	region             string
	compartment_id     string
	auth_token         string
	database_user_name string
	database_password  string
	//data_path             string
	//database_collections  string
	//database_tables       string
}

var (
	cfgFile        string
	cfgFileProfile string

	ociConfig ociConfigT
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "adb-cli",
	Short:   "adb-cli is a lightweight cli to manage Autonomous Databases in your OCI tenancy",
	Long:    "adb-cli is a lightweight cli to manage Autonomous Databases in your OCI tenancy",
	Version: cliVersion,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Print a banner
	banner := "ADB-Cli v" + cliVersion
	utils.PrintBold(banner)
	utils.Print(strings.Repeat("-", len(banner)))

	// cli entry point
	err := rootCmd.Execute()
	if err != nil {
		utils.PrintError("Error: " + err.Error())
		os.Exit(1)
	}

	// Generate docs in markdown language if requested by user
	if generateDocs, _ := rootCmd.Flags().GetBool("generate-docs"); generateDocs {
		path := "./docs"
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				utils.PrintError("Error: " + err.Error())
				os.Exit(1)
			}
		}
		rootCmd.DisableAutoGenTag = true
		err := doc.GenMarkdownTree(rootCmd, path)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			os.Exit(1)
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "adb-cli.yaml", "define the config file to use")
	rootCmd.PersistentFlags().StringVarP(&cfgFileProfile, "profile", "p", "DEFAULT", "define the profile to use")
	if runtime.GOOS == "windows" {
		// Disable the colours in the console for Windows OS family
		rootCmd.PersistentFlags().BoolVar(&utils.NoColor, "no-color", true, "disable color output")
	} else {
		rootCmd.PersistentFlags().BoolVar(&utils.NoColor, "no-color", false, "disable color output")
	}
	rootCmd.PersistentFlags().BoolVar(&utils.Verbose, "verbose", false, "increase verbosity")

	rootCmd.Flags().Bool("generate-docs", false, "generate docs in Markdown format")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag...
		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)
	} else {
		// ...or use the default one named adb-cli.config in the current directory.
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName("adb-cli")
	}

	if cfgFileProfile == "" {
		// If the profile is not specified use "DEFAULT" profile.
		cfgFileProfile = "DEFAULT"
	}

	// If a config file is found, read it.
	if err := viper.ReadInConfig(); err == nil {
		// These are printed only if the verbose flag is specified
		utils.PrintVerbose("Architecture          : " + runtime.GOARCH) // eventualmente spostarli altrove
		utils.PrintVerbose("OS family             : " + runtime.GOOS)
		utils.PrintVerbose("GO version            : " + runtime.Version())
		utils.PrintVerbose("OCI SDK for Go version: " + common.Version())
		utils.PrintVerbose("")

		utils.PrintInfo("Using config file: " + viper.ConfigFileUsed())
		utils.PrintInfo("Using profile    : " + cfgFileProfile)
		utils.PrintInfo("")

		// If there is no DEFAULT profile, print error and exit
		if !viper.GetViper().InConfig("DEFAULT") {
			utils.PrintError("The config file " + cfgFile + " does NOT contain the required DEFAULT profile")
			utils.PrintError("Plase add the default profile to the config file")
			utils.PrintError("To create a working template of the config file please run the following command after renaming or deleting the offending file:")
			utils.PrintError("adb-cli create config")
			os.Exit(1)
		}

		ociConfig.user = viper.GetViper().GetString(cfgFileProfile + ".user")
		ociConfig.pass_phrase = viper.GetViper().GetString(cfgFileProfile + ".pass_phrase")
		ociConfig.fingerprint = viper.GetViper().GetString(cfgFileProfile + ".fingerprint")
		ociConfig.tenancy = viper.GetViper().GetString(cfgFileProfile + ".tenancy")
		ociConfig.region = viper.GetViper().GetString(cfgFileProfile + ".region")
		ociConfig.compartment_id = viper.GetViper().GetString(cfgFileProfile + ".compartment_id")
		ociConfig.auth_token = viper.GetViper().GetString(cfgFileProfile + ".auth_token")
		ociConfig.database_user_name = viper.GetViper().GetString(cfgFileProfile + ".database_user_name")
		ociConfig.database_password = viper.GetViper().GetString(cfgFileProfile + ".database_password")
		if ociConfig.database_password != "" {
			if ociConfig.database_password[0] == '{' && ociConfig.database_password[len(ociConfig.database_password)-1] == '}' {
				ociConfig.database_password = utils.Decrypt(ociConfig.database_password[1 : len(ociConfig.database_password)-1])
			} else {
				viper.Set(cfgFileProfile+".database_password", "{"+utils.Encrypt(ociConfig.database_password)+"}")
				err = viper.WriteConfig()
				if err != nil {
					utils.Print("***** NOT DONE" + err.Error())
				}
			}

		}

		key_file, _ := homedir.Expand(viper.GetViper().GetString(cfgFileProfile + ".key_file"))
		key, err := os.ReadFile(key_file)
		if err != nil {
			// Exit if key file is not found
			utils.PrintError("Unable to read the private key file: " + err.Error())
			utils.PrintError("Please check the value of key_file in the config file: " + viper.ConfigFileUsed())
			os.Exit(1)
		}
		ociConfig.key_file = string(key)
	} else {
		// Exit if config file is not found
		utils.PrintError("Unable to read the config file: " + err.Error())
		utils.PrintError("Please create a config file template using the following command:")
		utils.PrintError("adb-cli create config [--with-keys]")
		utils.PrintError("Specify the --with-keys flag is you want to also create the keypair to be used with OCI")
		os.Exit(1)
	}
}
