package main

import (
	"net/url"
	"reflect"
	"testing"
)

/*
 -------------------------------------
|      Tests for getH1FromHTML       |
 -------------------------------------
*/

func TestGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLComplex(t *testing.T) {
	inputBody := `<html><body>
		<h1>Pecunia non olet</h1>
		<h2>Veritas<h2>
		<p>De veritas est maleficus</p>
	</body>
	</html>`
	actual := getH1FromHTML(inputBody)
	expected := "Pecunia non olet"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLNested(t *testing.T) {
	inputBody := `<html><body>
		<div>
			<h1>Header</h1>
		</div>
	</body></html>`
	actual := getH1FromHTML(inputBody)
	expected := "Header"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

/*
 -------------------------------------
| Tests for getFirstParagraphFromHTML |
 -------------------------------------
*/

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLWithoutMain(t *testing.T) {
	inputBody := `<html><body>
		<p>First Paragraph</p>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First Paragraph"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNestedMain(t *testing.T) {
	inputBody := `<html><body>
	<div>
		<main>
			<p>In vita umbratica</p>
		</main>
	</div>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "In vita umbratica"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

/*
 -------------------------------------
|     Tests for getURLsfromHTML      |
 -------------------------------------
*/

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLWithTwoLinks(t *testing.T) {
	inputURL := "https://www.google.com"
	inputBody := `<html>
<body>
<a href="/dev/avatar.png" target="_blank">Google</a>
<a href="/dev/logo.png" rel="noopener">GitHub</a>
</body>
</html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("coldn't parse input URL: %v", err)
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://www.google.com/dev/avatar.png",
		"https://www.google.com/dev/logo.png",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLWithThreeLinks(t *testing.T) {
	inputURL := "https://www.google.com"
	inputBody := `<html>
<body>
<a href="/pagina1.html" title="Primeira página">Página 1</a>
<a href="/pagina2.html" title="Primeira página">Página 1</a>
<a href="/pagina3.html" title="Primeira página">Página 1</a>
</body>
</html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://www.google.com/pagina1.html",
		"https://www.google.com/pagina2.html",
		"https://www.google.com/pagina3.html",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

/*
 -------------------------------------
|     Tests for getImagesFromHTML     |
 -------------------------------------
*/

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLTwoRelative(t *testing.T) {
	inputURL := "https://www.google.com"
	inputBody := `
	<html>
<body>
<img src="/assets/avatar.jpg" alt="Foto de perfil" width="100" height="100">
<img src="/assets/face/avatar.jpg" alt="Foto de perfil" width="100" height="100">

</body>
</html>
	`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://www.google.com/assets/avatar.jpg",
		"https://www.google.com/assets/face/avatar.jpg",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
