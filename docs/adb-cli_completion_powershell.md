## adb-cli completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	adb-cli completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
adb-cli completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string    define the config file to use (default "adb-cli.yaml")
      --no-color         disable color output
  -p, --profile string   define the profile to use (default "DEFAULT")
      --verbose          increase verbosity
```

### SEE ALSO

* [adb-cli completion](adb-cli_completion.md)	 - Generate the autocompletion script for the specified shell

