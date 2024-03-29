## adb-cli completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	adb-cli completion fish | source

To load completions for every new session, execute once:

	adb-cli completion fish > ~/.config/fish/completions/adb-cli.fish

You will need to start a new shell for this setup to take effect.


```
adb-cli completion fish [flags]
```

### Options

```
  -h, --help              help for fish
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

