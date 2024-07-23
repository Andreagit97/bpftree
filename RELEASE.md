# Release

## Versioning

- Bump a `MAJOR` (TBD).
- Bump a `MINOR` when you add at least one new field. In this way, old capture files become incompatible since we now support new fields.
- Bump a `PATCH` when you just need a fix inside the tool. For example a verifier issue or a bug in the userspace side.

## Release steps

1. update supported fields in the README.md with `bpftree f` output
2. update `release_changelog.md` file with `bpftree f` output
3. bump `bpftree` version in the code `toolVersion`
4. push a new tag (see next step)

## Push a tag

On the main branch

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
