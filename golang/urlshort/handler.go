package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			//http.Redirect(w,r,url, http.StatusSeeOther)
			http.Redirect(w, r, url, http.StatusFound) // don't use 301, 301 will keep the url in client
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	shortUrls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	m := buildMap(shortUrls)

	return MapHandler(m, fallback), nil
}

func parseYaml(yml []byte) ([]shortUrl, error) {
	shortUrls := make([]shortUrl, 0)

	if err := yaml.Unmarshal(yml, &shortUrls); err != nil {
		return nil, err
	}
	return shortUrls, nil
}

func buildMap(shortUrls []shortUrl) map[string]string {
	m := make(map[string]string)
	for _, url := range shortUrls {
		m[url.Path] = url.Url
	}
	return m
}

type shortUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
