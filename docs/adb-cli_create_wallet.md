## adb-cli create wallet

Create the wallet for an Autonomous Database identified by the --name flag

### Synopsis

Create the wallet for an Autonomous Database identified by the --name flag

```
adb-cli create wallet [flags]
```

### Options

```
  -h, --help              help for wallet
  -n, --name string       the name of the Autonomous Database for which create the wallet (required)
      --password string   the password that protects the wallet
```

### Options inherited from parent commands

```
  -c, --config string    define the config file to use (default "adb-cli.yaml")
      --no-color         disable color output
  -p, --profile string   define the profile to use (default "DEFAULT")
      --verbose          increase verbosity
```

### SEE ALSO

* [adb-cli create](adb-cli_create.md)	 - Create an OCI resource and/or supporting configuration - allowed resources: config, database, wallet

