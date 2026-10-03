# melodyme-ui

What MelodyMe's apps share. For now that is one page.

## conductor

MelodyMe's page for an address a visitor can't use: a site they may not read, or nothing at all. It is red, with a moving stave and a round conductor whose eyes follow the pointer (and cross, dizzily, when the pointer is between them), and it carries one message. The app shows it itself, so the status code (401, 403, 404) and the address stay its own.

```go
import "github.com/bobleung/melodyme-ui/conductor"

// Once, with the app's routes: the page's styles, script and fonts, from the app's own address.
mux.Handle("GET "+conductor.Prefix+"{path...}", conductor.Static())

// Nothing at the address: the generic page.
w.WriteHeader(http.StatusNotFound)
conductor.Page(conductor.NotFound).Render(ctx, w)

// A site the visitor may not read: say which, and offer sign-in when they aren't signed in.
conductor.Page(conductor.Stage{
	Title:    "DDSI",
	Subtitle: "Det Danske Suzuki Institut",
	Heading:  "Sorry, you don't have permission to visit this site",
	Text:     "Only people Det Danske Suzuki Institut has given access can visit.",
	SignIn:   &conductor.Form{Action: "/session/visit", Fields: []conductor.Field{{Name: "return_to", Value: "/ddsi/"}}},
})
```

The page needs nothing from the app's layout or styles, and no inline scripts or styles, so it works under a Content-Security-Policy of `'self'`.

## Working here

- `mise run check` is what every change must pass: generated templ files fresh, gofmt, go vet, tests.
- An app picks up a change only when it asks for the new version, so tag each release: `git tag v0.1.1 && git push --tags`, then `go get github.com/bobleung/melodyme-ui@v0.1.1` in the app.
- To try a change in an app before tagging it, point the app at this checkout with `go work init . ../melodyme-ui` there (and don't commit the `go.work`).
- The fonts are DM Serif Display, Fredoka and a subset of Noto Music, under the SIL Open Font License; their licences are in `conductor/static/fonts`.
