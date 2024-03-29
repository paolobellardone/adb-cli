## adb-cli delete

Delete an OCI resource and/or supporting configuration - allowed resources: config, database, wallet

### Synopsis

Delete an OCI resource and/or supporting configuration.
The allowed resources: config, database, wallet

```
adb-cli delete [flags]
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -c, --config string    define the config file to use (default "adb-cli.yaml")
      --no-color         disable color output
  -p, --profile string   define the profile to use (default "DEFAULT")
      --verbose          increase verbosity
```

### SEE ALSO

* [adb-cli](adb-cli.md)	 - adb-cli is a lightweight cli to manage Autonomous Databases in your OCI tenancy
* [adb-cli delete config](adb-cli_delete_config.md)	 - Delete the configuration file specified by --config global flag
* [adb-cli delete database](adb-cli_delete_database.md)	 - Delete an Autonomous Database identified by the --name flag
* [adb-cli delete wallet](adb-cli_delete_wallet.md)	 - Delete the wallet for an Autonomous Database identified by the --name flag

