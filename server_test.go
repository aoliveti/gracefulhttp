package gracefulhttp

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

type delayedHandler struct {
	http.Handler
	delay time.Duration
}

func (d *delayedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = r
	time.Sleep(d.delay)
	_, _ = w.Write([]byte("{}"))
}

func TestGracefulServer_ListenAndServeWithShutdown(t *testing.T) {
	t.Parallel()

	t.Run("gracefully shutdown on time", func(t *testing.T) {
		host := "localhost:34562"

		s := Bind(host, &delayedHandler{
			delay: 500 * time.Millisecond,
		})

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- s.ListenAndServeWithShutdown(ctx)
		}()

		r, err := http.Get("http://" + host)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		cancel()

		require.NoError(t, <-done)
	})

	t.Run("gracefully shutdown on time with cloudflare timeouts", func(t *testing.T) {
		host := "localhost:34563"

		s := Bind(host, &delayedHandler{
			delay: 500 * time.Millisecond,
		})

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- s.ListenAndServeWithShutdown(ctx, WithCloudflareTimeouts())
		}()

		r, err := http.Get("http://" + host)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		cancel()

		require.NoError(t, <-done)
	})

	t.Run("forcefully shutdown after a timeout", func(t *testing.T) {
		host := "localhost:34564"

		s := Bind(host, &delayedHandler{
			delay: 10 * time.Second,
		})

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error)
		go func() {
			done <- s.ListenAndServeWithShutdown(ctx)
		}()

		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

		_, err := http.Get("http://" + host)
		require.Error(t, err)

		require.ErrorIs(t, <-done, context.DeadlineExceeded)
	})
}

func TestGracefulServer_ListenAndServeTLSWithShutdown(t *testing.T) {
	const CertFile = "certs/cert.pem"
	const KeyFile = "certs/key.pem"

	if _, err := os.Stat(CertFile); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s not found", CertFile)
		return
	}

	if _, err := os.Stat(KeyFile); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s not found", KeyFile)
		return
	}

	t.Run("gracefully shutdown on time", func(t *testing.T) {
		host := "localhost:34572"

		s := Bind(host, &delayedHandler{
			delay: 500 * time.Millisecond,
		})

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- s.ListenAndServeTLSWithShutdown(ctx, CertFile, KeyFile)
		}()

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		r, err := http.NewRequest(http.MethodGet, "https://"+host, nil)
		require.NoError(t, err)

		response, err := client.Do(r)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		cancel()

		require.NoError(t, <-done)
	})

	t.Run("gracefully shutdown on time with cloudflare configurations", func(t *testing.T) {
		host := "localhost:34573"

		s := Bind(host, &delayedHandler{
			delay: 500 * time.Millisecond,
		})

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- s.ListenAndServeTLSWithShutdown(
				ctx,
				CertFile,
				KeyFile,
				WithCloudflareTimeouts(),
				WithCloudflareTLSConfig(),
			)
		}()

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		r, err := http.NewRequest(http.MethodGet, "https://"+host, nil)
		require.NoError(t, err)

		response, err := client.Do(r)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		cancel()

		require.NoError(t, <-done)
	})

	t.Run("forcefully shutdown a timeout", func(t *testing.T) {
		host := "localhost:34574"

		s := Bind(host, &delayedHandler{
			delay: 10 * time.Second,
		})

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error)
		go func() {
			done <- s.ListenAndServeTLSWithShutdown(ctx, CertFile, KeyFile)
		}()

		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		r, err := http.NewRequest(http.MethodGet, "https://"+host, nil)
		require.NoError(t, err)

		_, err = client.Do(r)
		require.Error(t, err)

		require.ErrorIs(t, <-done, context.DeadlineExceeded)
	})
}

func TestBind(t *testing.T) {
	handler := &delayedHandler{}
	type args struct {
		addr    string
		handler http.Handler
	}
	tests := []struct {
		name string
		args args
		want *GracefulServer
	}{
		{
			name: "bind server",
			args: args{
				addr:    "localhost:45678",
				handler: handler,
			},
			want: &GracefulServer{
				Server: http.Server{
					Addr:    "localhost:45678",
					Handler: handler,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Bind(tt.args.addr, tt.args.handler), "Bind(%v, %v)", tt.args.addr, tt.args.handler)
		})
	}
}
