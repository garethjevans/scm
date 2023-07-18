## scm token

Determines the token to use for an scm provider

```
scm token [flags]
```

### Examples

```
scm token --host=https://github.com --path .git-credentials
```

### Options

```
      --host string    The host of the scm provider, including scheme
  -k, --kind string    The kind of the scm provider
  -o, --owner string   The owner of the repository
  -p, --path string    The path to the git-credentials file
  -r, --repo string    The name of the repository
```

### Options inherited from parent commands

```
  -v, --debug   Debug Output
      --help    Show help for command
```

### SEE ALSO

* [scm](scm.md)	 - provides commands for interacting with different scm providers

