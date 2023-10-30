## scm monorepo-change

Determines the sha to clone for the supplied path on the monorepository

```
scm monorepo-change [flags]
```

### Examples

```
scm monorepo-change --host=https://github.com --path .git-credentials
```

### Options

```
      --branch string            The branch to search on (default: main) (default "main")
      --host string              The host of the scm provider, including scheme
  -k, --kind string              The kind of the scm provider
  -o, --owner string             The owner of the repository
  -p, --path string              The path to the git-credentials file
      --previous-commit string   The previous commit to search from (optional)
  -r, --repo string              The name of the repository
      --subpath string           The subPath to look for changes in (default: "")
```

### Options inherited from parent commands

```
  -v, --debug   Debug Output
      --help    Show help for command
```

### SEE ALSO

* [scm](scm.md)	 - provides commands for interacting with different scm providers

