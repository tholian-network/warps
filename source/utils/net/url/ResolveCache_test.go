package url

import net_url "net/url"
import "testing"

func TestResolveCache(t *testing.T) {

	t.Run("URL to host", func(t *testing.T) {

		url, _:= net_url.Parse("http://localhost:8080/")
		expected := "localhost:8080/index.html"
		resolved := ResolveCache(url)

		if resolved != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, resolved)
		}

	})

	t.Run("URL to folder", func(t *testing.T) {

		url, _:= net_url.Parse("http://localhost:8080/path/to/folder/")
		expected := "localhost:8080/path/to/folder/index.html"
		resolved := ResolveCache(url)

		if resolved != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, resolved)
		}

	})

	t.Run("URL to folder with parameters", func(t *testing.T) {

		url, _:= net_url.Parse("http://localhost:8080/path/to/folder/?a=b&c=d")
		expected := "localhost:8080/path/to/folder/index.html?a=b&c=d"
		resolved := ResolveCache(url)

		if resolved != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, resolved)
		}

	})

	t.Run("URL to folder with tracking parameters", func(t *testing.T) {

		url, _:= net_url.Parse("http://localhost:8080/path/to/folder/?utm_source=twitter&utm_campaign=1234")
		expected := "localhost:8080/path/to/folder/index.html"
		resolved := ResolveCache(url)

		if resolved != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, resolved)
		}

	})

	t.Run("URL to file", func(t *testing.T) {

		url, _:= net_url.Parse("http://localhost:8080/path/to/folder/and/some-file.jpg")
		expected := "localhost:8080/path/to/folder/and/some-file.jpg"
		resolved := ResolveCache(url)

		if resolved != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, resolved)
		}

	})

	t.Run("URL to file with parameters", func(t *testing.T) {

		url, _:= net_url.Parse("http://localhost:8080/path/to/folder/and/some-file.jpg?a=b&c=d")
		expected := "localhost:8080/path/to/folder/and/some-file.jpg?a=b&c=d"
		resolved := ResolveCache(url)

		if resolved != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, resolved)
		}

	})

	t.Run("URL to file with tracking parameters", func(t *testing.T) {

		url, _:= net_url.Parse("http://localhost:8080/path/to/folder/and/some-file.jpg?utm_source=twitter&utm_campaign=1234")
		expected := "localhost:8080/path/to/folder/and/some-file.jpg"
		resolved := ResolveCache(url)

		if resolved != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, resolved)
		}

	})

}
