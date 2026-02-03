# Contributing

Thanks for considering a contribution to github-kanri (gkn).

## Development

Requirements:

- Go 1.22+
- Git
- macOS

Run tests:

```sh
go test ./...
```

## Style

- `gofmt` required
- Short, explicit names
- Keep dependencies minimal
- Safety-first: deny rules must remain strict

## Issues

- Use GitHub Issues for bugs and feature requests
- For security reports, use `SECURITY.md` (no public issues)

## Pull requests

- Small, focused changes
- Include tests where possible (regressions preferred)
- Use Conventional Commits (`feat:`, `fix:`, `docs:`, ...)

## Safety bar

This project is designed for safe, non-interactive operations. Changes that weaken safety guarantees will not be accepted.
