# colly-linkcheck

Linkcheker that crawls ems.press and checks for dead links.

Right now the url and the ignore patterns are hardcoded.
Maybe we change that later so that can be customized.

## use directly

```sh
go run main.go > /dev/null # only print errors
go run main.go 2> >(tee errors.log) # save all errors to a file
```

## use as github action

```yaml
jobs:
    linkcheck:
        runs-on: ubuntu-latest
        name: check links
        steps:
        - name: Linkcheck
            # pin to current commit
            uses: ems-press/colly-linkcheck@main
```
