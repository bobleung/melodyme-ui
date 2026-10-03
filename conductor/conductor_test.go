package conductor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/bobleung/melodyme-ui/conductor"
)

func render(t *testing.T, st conductor.Stage) string {
	t.Helper()
	var b strings.Builder
	if err := conductor.Page(st).Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestTheGenericPageLinksBack(t *testing.T) {
	page := render(t, conductor.NotFound)
	for _, want := range []string{"<title>Wrong note</title>", "nothing at this address", `href="/"`, "powered by"} {
		if !strings.Contains(page, want) {
			t.Errorf("the not-found page lacks %q", want)
		}
	}
	if strings.Contains(page, "Sign in") {
		t.Error("the not-found page offers sign-in")
	}
}

func TestASiteCanOfferSignIn(t *testing.T) {
	page := render(t, conductor.Stage{
		Lang: "da", Title: "DDSI", Subtitle: "Det Danske Suzuki Institut",
		Heading: "Sorry, you don't have permission to visit this site", Text: "Only staff can visit.",
		SignIn: &conductor.Form{Action: "/session/visit", Fields: []conductor.Field{{Name: "return_to", Value: "/ddsi/"}}},
	})
	for _, want := range []string{`lang="da"`, "Det Danske Suzuki Institut", `action="/session/visit"`, `name="return_to" value="/ddsi/"`, "Sign in with MelodyMe ID"} {
		if !strings.Contains(page, want) {
			t.Errorf("the sign-in page lacks %q", want)
		}
	}
}

// Every file the page asks for is served, and the stylesheet's fonts too.
func TestTheFilesThePageNeedsAreServed(t *testing.T) {
	page := render(t, conductor.NotFound)
	srv := httptest.NewServer(conductor.Static())
	defer srv.Close()
	get := func(path string) *http.Response {
		res, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		return res
	}
	refs := regexp.MustCompile(`(?:href|src)="(`+regexp.QuoteMeta(conductor.Prefix)+`[^"]+)"`).FindAllStringSubmatch(page, -1)
	if len(refs) != 2 {
		t.Fatalf("want the stylesheet and the script, got %v", refs)
	}
	for _, m := range refs {
		path := strings.ReplaceAll(m[1], "&amp;", "&")
		if res := get(path); res.StatusCode != http.StatusOK || !strings.Contains(res.Header.Get("Cache-Control"), "immutable") {
			t.Errorf("%s: %d %q", path, res.StatusCode, res.Header.Get("Cache-Control"))
		}
	}
	for _, font := range []string{"dm-serif-display-latin.woff2", "fredoka-600-latin.woff2", "noto-music-notes.woff2"} {
		if res := get(conductor.Prefix + "fonts/" + font); res.StatusCode != http.StatusOK {
			t.Errorf("font %s: %d", font, res.StatusCode)
		}
	}
	if res := get(conductor.Prefix + "nope.css"); res.StatusCode != http.StatusNotFound {
		t.Errorf("an unknown file: %d", res.StatusCode)
	}
}
