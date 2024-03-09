package gracefulhttp

import (
	"crypto/tls"
	"time"
)

// GracefulServerOption is an option used to configure a [GracefulServer] instance.
type GracefulServerOption func(s *GracefulServer)

// WithShutdownTimeout sets the timeout for a graceful shutdown, after which all active connections
// will be forcibly closed.
func WithShutdownTimeout(duration time.Duration) GracefulServerOption {
	return func(s *GracefulServer) {
		if duration <= 0 {
			s.gracefulTimeout = defaultGracefulTimeout
			return
		}

		s.gracefulTimeout = duration
	}
}

// WithCloudflareTimeouts applies timeout patches to a [http.Server], implementing best practice
// configurations inspired by Cloudflare: https://blog.cloudflare.com/exposing-go-on-the-internet/
func WithCloudflareTimeouts() GracefulServerOption {
	return func(s *GracefulServer) {
		s.ReadTimeout = defaultReadTimeout
		s.ReadHeaderTimeout = defaultReadHeaderTimeout
		s.WriteTimeout = defaultWriteTimeout
		s.IdleTimeout = defaultIdleTimeout
	}
}

// WithCloudflareTLSConfig applies TLS configuration patches to a [http.Server], implementing best practice
// configurations inspired by Cloudflare: https://blog.cloudflare.com/exposing-go-on-the-internet/
func WithCloudflareTLSConfig() GracefulServerOption {
	return func(s *GracefulServer) {
		if s.TLSConfig == nil {
			s.TLSConfig = &tls.Config{}
		}

		s.TLSConfig.MinVersion = defaultTLSMinVersion
		s.TLSConfig.CurvePreferences = defaultTLSCurvePreferences
		s.TLSConfig.CipherSuites = defaultTLSCipherSuites
	}
}

// WithTLSConfig sets the provided TLS configuration for a [http.Server].
func WithTLSConfig(config *tls.Config) GracefulServerOption {
	return func(s *GracefulServer) {
		s.TLSConfig = config
	}
}
