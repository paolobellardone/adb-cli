## adb-cli create database

Create an Autonomous Database identified by the --name flag

### Synopsis

Create an Autonomous Database identified by the --name flag

```
adb-cli create database [flags]
```

### Options

```
  -m, --compute-model string         the compute model to be used to create the Autonomous Database - allowed values: ECPU, OCPU -- not used for Free Tier (default "ECPU")
  -u, --compute-units int            the number of compute units to allocate for the Autonomous Database -- not used for Free Tier (default 2)
      --enable-storage-autoscaling   enable autoscaling for storage (max 3x the size of reserved storage) -- not used for Free Tier
  -h, --help                         help for database
  -l, --license-model string         the licensing model to use - allowed values: full, byolee, byolse -- not used for Free Tier (default "full")
  -n, --name string                  the name of the Autonomous Database to create (required)
  -s, --storage int                  the size of storage in TB to allocate for the Autonomous Database -- not used for Free Tier (default 1)
  -t, --type string                  the type of the Autonomous Database to create - allowed values: atpfree, ajdfree, apexfree, adwfree, atp, ajd, apex, adw (default "atpfree")
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

