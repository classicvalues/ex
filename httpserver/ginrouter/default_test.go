package ginrouter

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"

	"github.com/circleci/ex/httpclient"
	"github.com/circleci/ex/httpserver"
	"github.com/circleci/ex/testing/testcontext"
)

func TestMiddleware(t *testing.T) {
	ctx, cancel := context.WithCancel(testcontext.Background())
	defer cancel()

	r := Default(ctx, "test server")
	r.GET("/foo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	srv, err := httpserver.New(ctx, "test-server", "127.0.0.1:0", r)
	assert.Assert(t, err)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return srv.Serve(ctx)
	})
	t.Cleanup(func() {
		assert.Check(t, g.Wait())
	})

	client := httpclient.New(httpclient.Config{
		Name:    "test-client",
		BaseURL: "http://" + srv.Addr(),
	})

	t.Run("Check we can get a 200 response", func(t *testing.T) {
		err = client.Call(ctx, httpclient.NewRequest("GET", "/foo", time.Second))
		assert.Assert(t, err)
	})
}
