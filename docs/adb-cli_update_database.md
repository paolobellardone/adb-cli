## adb-cli update database

Update the parameters of an Autonomous Database identified by the --name flag

### Synopsis

Update the parameters of an Autonomous Database identified by the --name flag

```
adb-cli update database [flags]
```

### Options

```
  -u, --compute-units int             the number of compute units to allocate for the Autonomous Database - not used for Free Tier (default 2)
      --disable-compute-autoscaling   disable autoscaling for compute units
      --disable-storage-autoscaling   disable autoscaling for storage
      --enable-compute-autoscaling    enable autoscaling for compute units (max 3x the number of baseline compute units)
      --enable-storage-autoscaling    enable autoscaling for storage (max 3x the size of reserved storage)
  -h, --help                          help for database
  -l, --license-model string          the licensing model to use - allowed values: full, byolee, byolse - not used for Free Tier (default "full")
  -n, --name string                   the name of the Autonomous Database to update (required)
  -s, --storage int                   the size of storage in GB to allocate for the Autonomous Database - not used for Free Tier (default 20)
  -t, --type string                   the type of the Autonomous Database to update - allowed values: atpfree, ajdfree, apexfree, adwfree, atp, ajd, apex, adw (default "atpfree")
```

### Options inherited from parent commands

```
  -c, --config string    define the config file to use (default "adb-cli.yaml")
      --no-color         disable color output
  -p, --profile string   define the profile to use (default "DEFAULT")
      --verbose          increase verbosity
```

### SEE ALSO

* [adb-cli update](adb-cli_update.md)	 - Update an OCI resource - allowed resources: database

