## scm tag

Creates a tag using the SCM api, will determine the repository url from the current directory, or can be overridden with --repository.

```
scm tag [flags]
```

### Examples

```
scm tag --tag 0.0.1 --sha abcdefgabcdefgabcdefgabcdefg --path .git-credentials
```

### Options

```
  -k, --kind string         The kind of the scm provider
  -p, --path string         The path to the git-credentials file
  -r, --repository string   The full https url of the repository
  -s, --sha string          The sha to attach the tag to
  -t, --tag string          The tag to apply
```

### Options inherited from parent commands

```
  -v, --debug   Debug Output
      --help    Show help for command
```

### SEE ALSO

* [scm](scm.md)	 - provides commands for interacting with different scm providers

