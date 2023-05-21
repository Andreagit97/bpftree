# Release

## Release steps

1. update supported fields in the README.md with `bpftree -f` output
2. update `release_changelog.md` file with `bpftree -f` output
3. bump `bpftree` version in the code
4. push a new tag

## Push a tag

```bash
git tag v0.0.1
git push origin v0.0.1
```

The CI should automatically craft a new release

## Delete a tag

```bash
git tag -d v0.0.1
git push origin --delete v0.0.1
```
