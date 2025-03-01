package httpserver

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	"github.com/circleci/ex/testing/testcontext"
)

func TestHTTPServer(t *testing.T) {
	ctx, cancel := context.WithCancel(testcontext.Background())
	defer cancel()

	r := http.NewServeMux()
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "hello world!")
	})

	srv, err := New(ctx, "test server", "localhost:0", r)
	assert.Assert(t, err)

	g, ctx := errgroup.WithContext(ctx)
	t.Cleanup(func() {
		assert.Check(t, g.Wait())
	})
	g.Go(func() error {
		return srv.Serve(ctx)
	})

	body, status := get(t, srv.Addr(), "test")
	assert.Check(t, cmp.Equal(status, http.StatusOK))
	assert.Check(t, cmp.Equal(body, "hello world!"))
}

func get(t *testing.T, baseurl, path string) (string, int) {
	t.Helper()

	r, err := http.Get(fmt.Sprintf("http://%s/%s", baseurl, path))
	assert.Assert(t, err)

	defer func() {
		assert.Assert(t, r.Body.Close())
	}()

	b, err := ioutil.ReadAll(r.Body)
	assert.Assert(t, err)

	return string(b), r.StatusCode
}
