package engine

import (
	"appengine"

	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func WriterCompact(w http.ResponseWriter, r *http.Request, c appengine.Context) io.Writer {
	// Don't work with GAE
	header := w.Header()
	var writer io.Writer = w
	c.Infof("Accept-Encoding: %s", r.Header.Get("Accept-Encoding"))
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		header.Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		writer = gz
		c.Infof("GZ")
	} else if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
		header.Set("Content-Encoding", "deflate")
		flate, err := flate.NewWriter(w, -1)
		if err == nil {
			defer flate.Close()
			writer = flate
			c.Infof("Flate")
		}
	} else {
		c.Infof("Uncompress")
	}
	return writer
}
