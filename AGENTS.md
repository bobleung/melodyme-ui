# melodyme-ui

What MelodyMe's apps (MelodyMe ID, Sites and School) share: a public Go module, `github.com/bobleung/melodyme-ui`. For now it is one package, `conductor`, the page an app shows for an address a visitor can't use. Read README.md first. AGENTS.md is the canonical agent instruction file; CLAUDE.md points here.

## Working here

- `mise run check` is what every change must pass: generated templ files fresh, gofmt, go vet, tests. `mise run gen` rebuilds the templ files, which are committed and never edited by hand.
- The repo is public. Nothing private goes in it: no secrets, no app's data, nothing specific to one app beyond what a `Stage` is given.
- Apps pick up a change only when they ask for a new version, so every change that apps should get is tagged (`v0.1.1`, `v0.2.0`, …) and then `go get` in each app. A change that breaks how apps call it (renaming a field, removing one) needs a minor version bump while below v1, and a note in the commit.
- The page must keep working under a Content-Security-Policy of `'self'`: no inline scripts, inline styles or `on*=` handlers, and nothing loaded from another host. Its styles, script and fonts are embedded and served by `conductor.Static()` at `conductor.Prefix`.
- The page renders the whole document itself, so it must not rely on any app's layout, `shell.css` or daisyUI.
- The app decides the status code and the wording. Text from a `Stage` is always escaped by templ; never render it as raw HTML.
- Standard library and templ only. Add a dependency only when the owner agrees.
- Fonts are under the SIL Open Font License; keep their licence files next to them in `conductor/static/fonts`.
