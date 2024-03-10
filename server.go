package gracefulhttp

import (
	"context"
	"crypto/tls"
	"errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

const (
	// defaultGracefulTimeout is the default timeout value used for a graceful shutdown.
	defaultGracefulTimeout = 5 * time.Second

	// defaultReadTimeout is the maximum duration for reading the entire request, including the body
	defaultReadTimeout = 5 * time.Second
	// defaultWriteTimeout is the maximum duration before timing out writes of the response
	defaultWriteTimeout = 10 * time.Second
	// defaultIdleTimeout is the maximum amount of time to wait for the next request when keepalive are enabled.
	defaultIdleTimeout = 120 * time.Second
	// defaultReadHeaderTimeout is the amount of time allowed to read request headers
	defaultReadHeaderTimeout = 5 * time.Second

	// defaultTLSMinVersion defines the recommended minimum version to use for the TLS protocol (1.2)
	defaultTLSMinVersion = tls.VersionTLS12
)

var (
	// defaultTLSCurvePreferences only use curves which have assembly implementations
	defaultTLSCurvePreferences = []tls.CurveID{
		tls.CurveP256,
		tls.X25519,
	}

	// defaultTLSCipherSuites defines the recommended cipher suites
	defaultTLSCipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

		// Best disabled, as they don't provide Forward Secrecy,
		// but might be necessary for some clients
		// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	}
)

// A GracefulServer is an extension of the [http.Server] that enables graceful shutdown.
// It allows the server to smoothly stop accepting new connections while
// processing existing requests to completion within a specified timeout.
// After the timeout, ongoing connections will be forcibly closed.
// The default timeout is set to 5 seconds
type GracefulServer struct {
	http.Server

	gracefulTimeout time.Duration
}

// Bind returns a new [GracefulServer] configured with the provided address and handler.
func Bind(addr string, handler http.Handler) *GracefulServer {
	return &GracefulServer{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// ListenAndServeWithShutdown starts a [http.Server] with the given address and handler.
// It blocks until the context is canceled or an error occurs.
// If the context is canceled, the server will attempt a graceful shutdown.
// If the graceful shutdown exceeds the provided timeout, the server will be forcefully closed.
// The default timeout is set to 5 seconds.
// The [context.Canceled] error is intentionally ignored and thus not returned by the method.
// Upon timeout, the method returns only [context.DeadlineExceeded] error.
func (s *GracefulServer) ListenAndServeWithShutdown(ctx context.Context, opts ...GracefulServerOption) error {
	s.initialize(opts)

	return s.listenAndServe(ctx)
}

// ListenAndServeTLSWithShutdown starts a [http.Server] with the provided address, handler, certificate, and key.
// It behaves similarly to [ListenAndServeWithShutdown] but for HTTPS connections.
// For additional details, refer to the documentation of [ListenAndServeWithShutdown].
func (s *GracefulServer) ListenAndServeTLSWithShutdown(ctx context.Context, certFile string, keyFile string, opts ...GracefulServerOption) error {
	s.initialize(opts)

	return s.listenAndServeTLS(ctx, certFile, keyFile)
}

// listenAndServe invokes [http.ListenAndServe] until the context is canceled, then invokes the shutdown method.
func (s *GracefulServer) listenAndServe(ctx context.Context) error {
	g := errgroup.Group{}

	g.Go(func() error {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})
	g.Go(func() error {
		<-ctx.Done()

		return s.shutdown()
	})

	return g.Wait()
}

// listenAndServeTLS invokes [http.ListenAndServeTLS] with the supplied certificate and key files until the context is canceled,
// then invokes the shutdown method.
func (s *GracefulServer) listenAndServeTLS(ctx context.Context, certFile string, keyFile string) error {
	g := errgroup.Group{}

	g.Go(func() error {
		if err := s.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})
	g.Go(func() error {
		<-ctx.Done()

		return s.shutdown()
	})

	return g.Wait()
}

// initialize set the default timeout to 5s and sets the GracefulServer options
func (s *GracefulServer) initialize(opts []GracefulServerOption) {
	s.gracefulTimeout = defaultGracefulTimeout

	for _, opt := range opts {
		opt(s)
	}
}

// shutdown invokes [http.Shutdown], and if there is a timeout,
// it will forcibly close the active connections using [http.Close].
func (s *GracefulServer) shutdown() error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), s.gracefulTimeout)
	defer cancel()

	done := make(chan struct{}, 1)

	g, groupCtx := errgroup.WithContext(ctxTimeout)
	g.Go(func() error {
		defer close(done)
		return s.Shutdown(groupCtx)
	})
	g.Go(func() error {
		select {
		case <-groupCtx.Done():
			return s.Close()
		case <-done:
			return nil
		}
	})

	return g.Wait()
}
