## adb-cli completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(adb-cli completion bash)

To load completions for every new session, execute once:

#### Linux:

	adb-cli completion bash > /etc/bash_completion.d/adb-cli

#### macOS:

	adb-cli completion bash > $(brew --prefix)/etc/bash_completion.d/adb-cli

You will need to start a new shell for this setup to take effect.


```
adb-cli completion bash
```

### Options

```
  -h, --help              help for bash
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

