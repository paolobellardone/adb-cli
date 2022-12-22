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
package utils

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
)

// GetAutonomousDatabaseOCID: returns the OCID of an ADB
func GetAutonomousDatabaseOCID(client database.DatabaseClient, adbName string, compartmentOCID string) (string, error) {
	// TODO: use dbName instead of display name because it can be multiple ???
	// TODO: the ListAutonomousDatabase works on DisplayName to filter, if there are more than one list those and exit or make the user choose
	dbRequest := database.ListAutonomousDatabasesRequest{
		CompartmentId: common.String(compartmentOCID),
		DisplayName:   common.String(adbName),
	}

	dbResponse, err := client.ListAutonomousDatabases(context.Background(), dbRequest)
	if err != nil {
		return "", err
	}

	var adbOCID string
	if len(dbResponse.Items) < 1 {
		return "", errors.New("The Autonomous Database " + adbName + " was not found")
	} else if len(dbResponse.Items) > 1 {
		return "", errors.New("There are too many autonomous databases with the same name, it is not possible")
	} else {
		adbOCID = *dbResponse.Items[0].Id
		PrintVerbose("The Autonomous Database " + adbName + " has id: " + adbOCID)
		return adbOCID, nil
	}
}

// GetAutonomousDatabase: returns the details of an ADB
func GetAutonomousDatabase(client database.DatabaseClient, adbOCID string) (database.AutonomousDatabase, error) {
	dbRequest := database.GetAutonomousDatabaseRequest{
		AutonomousDatabaseId: common.String(adbOCID),
	}

	dbResponse, err := client.GetAutonomousDatabase(context.Background(), dbRequest)
	if err != nil {
		return dbResponse.AutonomousDatabase, err
	} else {
		return dbResponse.AutonomousDatabase, nil
	}
}

// DecodeADBTypeForCreate: converts the database type parameter to an Enum
func DecodeADBTypeForCreate(adbType string) (database.CreateAutonomousDatabaseBaseDbWorkloadEnum, bool) {
	/*
		atpfree : Always Free Autonomous Transaction Processing (default)
		ajdfree : Always Free Autonomous JSON Database
		apexfree: Always Free Autonomous Application Express
		adwfree : Always Free Autonomous Data Warehouse
		atp     : Autonomous Transaction Processing
		ajd     : Autonomous JSON Database
		apex    : Autonomous Application Express
		adw     : Autonomous Data Warehouse
	*/
	switch adbType {
	case "atpfree":
		return database.CreateAutonomousDatabaseBaseDbWorkloadOltp, true
	case "ajdfree":
		return database.CreateAutonomousDatabaseBaseDbWorkloadAjd, true
	case "apexfree":
		return database.CreateAutonomousDatabaseBaseDbWorkloadApex, true
	case "adwfree":
		return database.CreateAutonomousDatabaseBaseDbWorkloadDw, true
	case "atp":
		return database.CreateAutonomousDatabaseBaseDbWorkloadOltp, false
	case "ajd":
		return database.CreateAutonomousDatabaseBaseDbWorkloadAjd, false
	case "apex":
		return database.CreateAutonomousDatabaseBaseDbWorkloadApex, false
	case "adw":
		return database.CreateAutonomousDatabaseBaseDbWorkloadDw, false
	default:
		// The default is atpfree
		return database.CreateAutonomousDatabaseBaseDbWorkloadOltp, true
	}
}

// DecodeADBTypeForUpdate: converts the database type parameter to an Enum
func DecodeADBTypeForUpdate(adbType string) (database.UpdateAutonomousDatabaseDetailsDbWorkloadEnum, bool) {
	/*
		atpfree : Always Free Autonomous Transaction Processing (default)
		ajdfree : Always Free Autonomous JSON Database
		apexfree: Always Free Autonomous Application Express
		adwfree : Always Free Autonomous Data Warehouse
		atp     : Autonomous Transaction Processing
		ajd     : Autonomous JSON Database
		apex    : Autonomous Application Express
		adw     : Autonomous Data Warehouse
	*/
	switch adbType {
	case "atpfree":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadOltp, true
	case "ajdfree":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadAjd, true
	case "apexfree":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadApex, true
	case "adwfree":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadDw, true
	case "atp":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadOltp, false
	case "ajd":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadAjd, false
	case "apex":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadApex, false
	case "adw":
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadDw, false
	default:
		// The default is atpfree
		return database.UpdateAutonomousDatabaseDetailsDbWorkloadOltp, true
	}
}

// DecodeLicenseModelForCreate: converts the database license model parameter to an Enum
func DecodeLicenseModelForCreate(license_model string) database.CreateAutonomousDatabaseBaseLicenseModelEnum {
	/*
	   full: License included
	   byol: Bring your own license
	*/
	switch license_model {
	case "full":
		return database.CreateAutonomousDatabaseBaseLicenseModelLicenseIncluded
	case "byol":
		return database.CreateAutonomousDatabaseBaseLicenseModelBringYourOwnLicense
	default:
		// The default is full
		return database.CreateAutonomousDatabaseBaseLicenseModelLicenseIncluded
	}
}

// DecodeLicenseModelForUpdate: converts the database license model parameter to an Enum
func DecodeLicenseModelForUpdate(license_model string) database.UpdateAutonomousDatabaseDetailsLicenseModelEnum {
	/*
	   full: License included
	   byol: Bring your own license
	*/
	switch license_model {
	case "full":
		return database.UpdateAutonomousDatabaseDetailsLicenseModelLicenseIncluded
	case "byol":
		return database.UpdateAutonomousDatabaseDetailsLicenseModelBringYourOwnLicense
	default:
		// The default is full
		return database.UpdateAutonomousDatabaseDetailsLicenseModelLicenseIncluded
	}
}

// CreateAutonomousDatabaseWallet: creates an ADB wallet file
func CreateAutonomousDatabaseWallet(client database.DatabaseClient, adbName string, adbOCID string, password string) (string, error) {
	generateADBWalletRequest := database.GenerateAutonomousDatabaseWalletRequest{
		AutonomousDatabaseId: common.String(adbOCID),
		GenerateAutonomousDatabaseWalletDetails: database.GenerateAutonomousDatabaseWalletDetails{
			GenerateType: database.GenerateAutonomousDatabaseWalletDetailsGenerateTypeSingle,
			Password:     common.String(password),
		},
	}

	generateADBWalletResponse, err := client.GenerateAutonomousDatabaseWallet(context.Background(), generateADBWalletRequest)
	if err != nil {
		return "", err
	}

	walletName := "Wallet_" + adbName + ".zip"
	adbWallet, err := os.Create(walletName)
	if err != nil {
		return "", err
	}
	defer adbWallet.Close()
	written, err := io.Copy(adbWallet, generateADBWalletResponse.Content)
	if err != nil || written != *generateADBWalletResponse.ContentLength {
		return "", err
	}

	return walletName, nil
}

// DeleteAutonmousDatabaseWallet: deletes an ADB wallet file
func DeleteAutonmousDatabaseWallet(adbName string) error {
	err := os.Remove("Wallet_" + adbName + ".zip")
	return err
}
