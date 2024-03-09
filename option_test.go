package gracefulhttp

import (
	"crypto/tls"
	"reflect"
	"testing"
	"time"
)

func TestWithCloudflareTimeouts(t *testing.T) {
	type timeouts struct {
		ReadTimeout       time.Duration
		ReadHeaderTimeout time.Duration
		WriteTimeout      time.Duration
		IdleTimeout       time.Duration
	}
	tests := []struct {
		name string
		want timeouts
	}{
		{
			name: "set cloudflare timeouts",
			want: timeouts{
				ReadTimeout:       defaultReadTimeout,
				ReadHeaderTimeout: defaultReadHeaderTimeout,
				WriteTimeout:      defaultWriteTimeout,
				IdleTimeout:       defaultIdleTimeout,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GracefulServer{}
			opt := WithCloudflareTimeouts()
			opt(&s)

			got := timeouts{
				ReadTimeout:       s.ReadTimeout,
				ReadHeaderTimeout: s.ReadHeaderTimeout,
				WriteTimeout:      s.WriteTimeout,
				IdleTimeout:       s.IdleTimeout,
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCloudflareTimeouts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithShutdownTimeout(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "set timeout to a positive value",
			args: args{
				duration: 7 * time.Second,
			},
			want: 7 * time.Second,
		},
		{
			name: "set timeout to zero",
			args: args{
				duration: 0 * time.Second,
			},
			want: defaultGracefulTimeout,
		},
		{
			name: "set a negative timeout",
			args: args{
				duration: -5 * time.Second,
			},
			want: defaultGracefulTimeout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GracefulServer{}
			opt := WithShutdownTimeout(tt.args.duration)
			opt(&s)

			if got := s.gracefulTimeout; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithShutdownTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCloudflareTLSConfig(t *testing.T) {
	type config struct {
		MinVersion       uint16
		CurvePreferences []tls.CurveID
		CipherSuites     []uint16
	}
	tests := []struct {
		name string
		want config
	}{
		{
			name: "set cloudflare tls config",
			want: config{
				MinVersion:       defaultTLSMinVersion,
				CurvePreferences: defaultTLSCurvePreferences,
				CipherSuites:     defaultTLSCipherSuites,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GracefulServer{}
			opt := WithCloudflareTLSConfig()
			opt(&s)

			got := config{
				MinVersion:       s.TLSConfig.MinVersion,
				CurvePreferences: s.TLSConfig.CurvePreferences,
				CipherSuites:     s.TLSConfig.CipherSuites,
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCloudflareTimeouts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLSConfig(t *testing.T) {
	type args struct {
		config *tls.Config
	}
	tests := []struct {
		name string
		args args
		want *tls.Config
	}{
		{
			name: "pass tls config",
			args: args{
				config: &tls.Config{
					MinVersion:       defaultTLSMinVersion,
					CurvePreferences: defaultTLSCurvePreferences,
					CipherSuites:     defaultTLSCipherSuites,
				},
			},
			want: &tls.Config{
				MinVersion:       defaultTLSMinVersion,
				CurvePreferences: defaultTLSCurvePreferences,
				CipherSuites:     defaultTLSCipherSuites,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GracefulServer{}
			opt := WithTLSConfig(tt.args.config)
			opt(&s)

			got := s.TLSConfig

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCloudflareTimeouts() = %v, want %v", got, tt.want)
			}
		})
	}
}
