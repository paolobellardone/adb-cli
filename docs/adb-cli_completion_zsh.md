## adb-cli completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(adb-cli completion zsh)

To load completions for every new session, execute once:

#### Linux:

	adb-cli completion zsh > "${fpath[1]}/_adb-cli"

#### macOS:

	adb-cli completion zsh > $(brew --prefix)/share/zsh/site-functions/_adb-cli

You will need to start a new shell for this setup to take effect.


```
adb-cli completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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

