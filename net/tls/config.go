package tls

import "time"

type TLSConfig struct {
	CertFile   string      `yaml:"certFile,omitempty" json:"certFile,omitempty"`
	KeyFile    string      `yaml:"keyFile,omitempty" json:"keyFile,omitempty"`
	CAFile     string      `yaml:"caFile,omitempty" json:"caFile,omitempty"`
	Secure     bool        `yaml:",omitempty" json:"secure,omitempty"`
	ServerName string      `yaml:"serverName,omitempty" json:"serverName,omitempty"`
	Options    *TLSOptions `yaml:",omitempty" json:"options,omitempty"`

	// for auto-generated default certificate.
	Validity     time.Duration `yaml:",omitempty" json:"validity,omitempty"`
	CommonName   string        `yaml:"commonName,omitempty" json:"commonName,omitempty"`
	Organization string        `yaml:",omitempty" json:"organization,omitempty"`
}

type TLSOptions struct {
	MinVersion   string   `yaml:"minVersion,omitempty" json:"minVersion,omitempty"`
	MaxVersion   string   `yaml:"maxVersion,omitempty" json:"maxVersion,omitempty"`
	CipherSuites []string `yaml:"cipherSuites,omitempty" json:"cipherSuites,omitempty"`
	ALPN         []string `yaml:"alpn,omitempty" json:"alpn,omitempty"`
}
