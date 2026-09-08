# Security policy

Atesaki is a v0 pre-release. Security fixes target `main` only. Older commits
and branches receive no backported fixes. The supported platforms are Linux
and macOS.

## Report a vulnerability

Use [GitHub private vulnerability reporting](https://github.com/acartag7/atesaki/security/advisories/new).
Do not put vulnerability details or credentials in public issues or commits.

Private vulnerability reporting was enabled with the repository owner's
authorization and verified on 2026-09-08.

The maintainer aims to acknowledge a private report within three business days
and provide an initial assessment within ten business days. Resolution timing
depends on severity and the available fix; these are response targets, not a
guarantee of a release date.

Reports about the binary, published container, and deployment recipe are in
scope. Include the affected commit or release, platform, expected behavior,
observed behavior, and a minimal reproduction with sensitive values removed.
The container and recipe remain planned until published.

## Coordinated fixes

Discuss a report and prepare its fix privately. Public commits must not narrate
an unreleased vulnerability before the fixed release is available. Coordinate
disclosure with the reporter and publish an advisory with the affected versions
and mitigation when the fix is released.
