## adb-cli create config

Create the configuration template file with name adb-cli.yaml.template

### Synopsis

Create the configuration template file with name adb-cli.yaml.template

Duplicate the file and edit as needed (documentation included) to configure the cli
to operate on your OCI tenancy.

```
adb-cli create config [flags]
```

### Options

```
  -h, --help        help for config
  -k, --with-keys   create key pair to be used with this configuration
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

