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
	"os"
	"strings"

	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"path/filepath"

	"github.com/paolobellardone/adb-cli/utils"

	"github.com/spf13/cobra"
)

// TODO: modify the template in order to adapt to the new scope of the cli

// This is the template of the configuration file needed by the cli
// Modify as needed in order to enrich the features of the cli
var template string = `# Copyright © 2022 PaoloB
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

# DEFAULT profile, you can define others, for example: ASHBURN_REGION or TEST_ENVIRONMENT
# You can choose a profile using the --profile command line argument
# This configuration file MUST have at least one profile named DEFAULT
default:

  # OCID of your tenancy. To get the value, see:
  # https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five
  tenancy: ocid1.tenancy.oc1..<unique_ID>

  # An Oracle Cloud Infrastructure region identifier. For a list of possible region identifiers, see:
  # https://docs.cloud.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm#top
  region=: eu-frankfurt-1

  # OCID of the compartment to use for resources creation. To get more information about compartments, see:
  # https://docs.cloud.oracle.com/en-us/iaas/Content/Identity/Tasks/managingcompartments.htm?Highlight=compartment%20ocid#Managing_Compartments
  compartment_id: ocid1.compartment.oc1..<unique_ID>

  # OCID of the user connecting to Oracle Cloud Infrastructure APIs. To get the value, see:
  # https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#five
  user: ocid1.user.oc1..<unique_ID>

  # Full path and filename of the SSH private key (use *solely* forward slashes).
  # /!\ Warning: The key pair must be in PEM format (2048 bits). For instructions on generating a key pair in PEM format, see:
  # https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#Required_Keys_and_OCIDs
  key_file: <full path to SSH private key file>

  # Fingerprint for the SSH *public* key that was added to the user mentioned above. To get the value, see:
  # https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm#four
  fingerprint: <fingerprint associated with the corresponding SSH *public* key>

  # Uncomment in the case your SSH private key needs a pass phrase.
  # The adb-cli generated keys does not need it.
  # pass_phrase: <pass phrase to use with your SSH private key>

  # Authentication token that will be used for OCI Object Storage configuration, see:
  # https://docs.cloud.oracle.com/en-us/iaas/Content/Registry/Tasks/registrygettingauthtoken.htm?Highlight=user%20auth%20tokens
  auth_token: <authentication token>

  # The database password used for database creation and admin user
  # - 12 chars minimum and 30 chars maximum
  # - can't contain the database user name word
  # - contains 1 digit minimum
  # - contains 1 lower case char
  # - contains 1 upper case char
  database_password: <database password>

  # Path to a folder where data to load into collections can be found (default to current directory)
  data_path: .

  # A list of comma separated JSON collection name(s) that you wish to get right after database creation
  # database_collections:

  # A list of comma separated table name(s) that you wish to get right after database creation
  # These table must have corresponding CSV file(s) so that table structure (DDL) is deduced from the files
  # database_tables=:
`

// create_configCmd creates the configuration template file for this cli
var create_configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create the configuration template file with name adb-cli.yaml.template",
	Long: `Create the configuration template file with name adb-cli.yaml.template

Duplicate the file and edit as needed (documentation included) to configure the cli
to operate on your OCI tenancy.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintVerbose("create config command called")
		utils.PrintVerbose("")

		withKeys, _ := cmd.Flags().GetBool("with-keys")
		if withKeys {
			utils.PrintInfo("Create configuration template with key pair for OCI")
			// TODO: print information on what to do with the keys in OCI
			private_key_path, public_key_fingerprint, err := create_keypair()
			if err != nil {
				utils.PrintError("Error: " + err.Error())
				return
			} else {
				template = strings.Replace(template, "key_file: <full path to SSH private key file>", "key_file: "+private_key_path, 1)
				template = strings.Replace(template, "fingerprint: <fingerprint associated with the corresponding SSH *public* key>", "fingerprint: "+public_key_fingerprint, 1)
			}
		} else {
			utils.PrintInfo("Create configuration template")
		}

		err := os.WriteFile("adb-cli.yaml.template", []byte(template), 0664)
		if err != nil {
			utils.PrintError("Error: " + err.Error())
			return
		} else {
			utils.PrintInfo("Done")
		}
	},
}

func init() {
	createCmd.AddCommand(create_configCmd)

	create_configCmd.Flags().BoolP("with-keys", "k", false, "create key pair to be used with this configuration")
}

// create_keypair: creates a key pair that can be used in OCI to sign the API calls
func create_keypair() (private_key_path string, public_key_fingerprint string, err error) {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		//fmt.Printf("Cannot generate RSA key\n")
		return
	}
	publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create("oci_user_private.pem")
	if err != nil {
		//fmt.Printf("error when create private.pem: %s \n", err)
		return
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		//fmt.Printf("error when encode private pem: %s \n", err)
		return
	}
	private_key_path, err = filepath.Abs("oci_user_private.pem")
	if err != nil {
		//fmt.Println("Error:", err)
		return
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		//fmt.Printf("error when dumping publickey: %s \n", err)
		return
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create("oci_user_public.pem")
	if err != nil {
		//fmt.Printf("error when create public.pem: %s \n", err)
		return
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		//fmt.Printf("error when encode public pem: %s \n", err)
		return
	}

	// generate fingerprint for public key and dump to file
	md5sum := md5.Sum(publicKeyBytes)
	hexarray := make([]string, len(md5sum))
	for i, c := range md5sum {
		hexarray[i] = hex.EncodeToString([]byte{c})
	}
	public_key_fingerprint = strings.Join(hexarray, ":")
	err = os.WriteFile("oci_user_public.fingerprint", []byte(public_key_fingerprint), 0664)
	if err != nil {
		return
	}

	// If not returned before due to an error, here returns the private key path and public key fingerprint (nil for error)
	return
}
