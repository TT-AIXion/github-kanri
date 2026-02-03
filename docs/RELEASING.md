# Releasing

## Preconditions

- Clean working tree
- `go test ./...` passes
- Optional: `golangci-lint run` and `govulncheck ./...`
- `CHANGELOG.md` updated

## Release steps (semver)

1. Update `Formula/gkn.rb`:
   - `version`
   - `tag` (`vX.Y.Z`)
   - `revision` (release commit SHA)
2. Commit release prep (`chore: release vX.Y.Z`).
3. Tag and push:

```sh
git tag -a vX.Y.Z -m "release vX.Y.Z"
git push origin vX.Y.Z
```

4. `release` workflow runs GoReleaser and publishes artifacts.
5. Verify release assets, checksums, SBOMs, and signatures (cosign bundle).

## Snapshot tags

`auto-tag` creates snapshot tags `v0.0.0-main.*` on `main` for CI builds. Do not use these for public releases.
