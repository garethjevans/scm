## scm pr create

Commits and pushes any local changes, if the base & commit branches are different, then a PR/MR will be created.

```
scm pr create [flags]
```

### Examples

```
scm pr create --commit-branch update .git-credentials
```

### Options

```
      --base-branch string        The branch to target the PR to (default "main")
      --commit-branch string      The branch to push the commits to
      --commit-title string       The title of any commits to push
      --git-email string          The author of any git commits
      --git-user string           The author of any git commits
      --kind string               In case we are unable to determine the type of scm server, this can provide hints
      --output-git-sha string     The location to write the git sha to
      --output-pr-number string   The location to write the pr number to
      --output-pr-url string      The location to write the pr url to
  -p, --path string               The path to the git-credentials file
      --pr-title string           The title of the PR/MR to create
```

### Options inherited from parent commands

```
  -v, --debug   Debug Output
      --help    Show help for command
```

### SEE ALSO

* [scm pr](scm_pr.md)	 - 

