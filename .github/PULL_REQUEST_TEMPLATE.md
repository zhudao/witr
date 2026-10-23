## Description

Briefly describe what this PR does. Mention the issue it fixes, if any (e.g. `Fixes #123`).

## Type of change

Check all that apply:

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (may affect existing functionality)
- [ ] This change requires a documentation update

## Checklist

- [ ] I have opened this PR against the `staging` branch, with a Conventional Commits title
- [ ] I have performed a self-review of my own code
- [ ] `make lint` and `make test` pass, and my changes add no new compiler, `go vet` or lint warnings
- [ ] Platform-specific changes are checked in every platform's file (`_linux`, `_darwin`, `_freebsd`, `_windows`)
- [ ] Features that exist in both the CLI and the TUI are updated in both
- [ ] Output changes are mirrored in the playground (`docs/js/engine.js`, fixtures regenerated)
- [ ] New `--json` fields are recorded in the JSON contract golden file
- [ ] I have made corresponding changes to the documentation (README, agent skill)

## Thank you for your contribution! 🎉

We’re excited to review your pull request. Please fill out the details above to help us understand your changes. Don’t worry if you can’t check every box, just do your best!
