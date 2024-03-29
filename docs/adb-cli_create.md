## adb-cli create

Create an OCI resource and/or supporting configuration - allowed resources: config, database, wallet

### Synopsis

Create an OCI resource and/or supporting configuration.
The allowed resources are: config, database, wallet

```
adb-cli create [flags]
```

### Options

```
  -h, --help   help for create
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
* [adb-cli create config](adb-cli_create_config.md)	 - Create the configuration template file with name adb-cli.yaml.template
* [adb-cli create database](adb-cli_create_database.md)	 - Create an Autonomous Database identified by the --name flag
* [adb-cli create wallet](adb-cli_create_wallet.md)	 - Create the wallet for an Autonomous Database identified by the --name flag

