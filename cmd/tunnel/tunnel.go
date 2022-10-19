// Copyright (C) 2017 Michał Matczuk
// Use of this source code is governed by an AGPL-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/cenkalti/backoff"
	tunnel "github.com/mmatczuk/go-http-tunnel"
	"github.com/mmatczuk/go-http-tunnel/id"
	"github.com/mmatczuk/go-http-tunnel/log"
	"github.com/mmatczuk/go-http-tunnel/proto"
)

var crt, key = `
-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUIEmCIcud0MupndOLGVHdiKVSEc0wDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjA3MjAwODM3MjZaFw0yMjA4
MTkwODM3MjZaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDUhdAjAruIo+nvm6FQMmviax1HbV8hhdxZBejBfwTr
Iw0se3uXAiP+64wUyG5OAh5+5jDSEGOkNUZEQuHBsFBv4UEh5T6vGUcIgQQSgIy4
2VMmirwe6Xj1R762+dezRb3L6JJmINOgJiSz05pBlxFqmKUlo5KVrDgDsPHUeBMC
sTNDw9x9RocjftD09OdHap1g9W06erequQWTVkFy0f5SXTU1GHCJ/vf02CmvkBIQ
eFG2K48NPeZIeKBQrbwQ5rW7yo+juUPqGyk1wjZIVJVT2xABIXvxZVKJ3txsDYcF
haWyiHLain99tidCPpuek21sxomU4JbojP1E9bJw/YOlAgMBAAGjUzBRMB0GA1Ud
DgQWBBSg61RIz23aRMDVnJwKZyjbQWTSgDAfBgNVHSMEGDAWgBSg61RIz23aRMDV
nJwKZyjbQWTSgDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQC0
7H4NA37fn3qKKFM0J9XTiJglt73SMsY6sYfNPdXQ60H5JvT8MrNdlpKTShBmhHbO
Mv6ASE9t7bHo5DejiypwaO+KjjKkAiQfUdkh7RrnP5y4rYWmgjuY64vEDqrNM8IF
4OkDmCiSlTwFVUHgCG0vtoJoCKJLRCv0fxXO7J3uRATdw2IvbItfoFggR1LWixrI
QrH7010WD4Id6DA67RBwSN3shTHjZ76njbWSaT3B4bdqR0+AJRBugw5+w9ipbwVI
yEvOIc4mmFcVqRr/FPmxURWUXOmjzuTGKHZzGgkStO/OI9RsV7fsUp6Tgnpn+/sQ
VLmG4tNC1u3I8LziPXm1
-----END CERTIFICATE-----`, `
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDUhdAjAruIo+nv
m6FQMmviax1HbV8hhdxZBejBfwTrIw0se3uXAiP+64wUyG5OAh5+5jDSEGOkNUZE
QuHBsFBv4UEh5T6vGUcIgQQSgIy42VMmirwe6Xj1R762+dezRb3L6JJmINOgJiSz
05pBlxFqmKUlo5KVrDgDsPHUeBMCsTNDw9x9RocjftD09OdHap1g9W06erequQWT
VkFy0f5SXTU1GHCJ/vf02CmvkBIQeFG2K48NPeZIeKBQrbwQ5rW7yo+juUPqGyk1
wjZIVJVT2xABIXvxZVKJ3txsDYcFhaWyiHLain99tidCPpuek21sxomU4JbojP1E
9bJw/YOlAgMBAAECggEBAMCtRQPI4W9DdL+IsNq9q/QOOvBcJ4iEDT8kzV5Io5Pv
Jm1y0p0ZErt2oHzjLqZI448XcaGxvDTPHBKdNIJfML4OUffCGG+1kiISnjeMFoks
d2eVFnNMJx7x2eWYZEgjbazmAXzYPCWRD2t95/eZV+/1zpxuTWKnTe2Cye1go5Om
ntwILe0xEVexPeqt+eunf4e7O7sZhZY3i311vT8ZKLI9gTeI0tr1iSIPIU1Wklep
WDSM8Ls/eFSaIWOT+uecTb3InllMnlNcZvZX6x0Enl+vIwj5QQv4e1YD5k5YzskD
9pZkWjTw/bwMJ1hZt5vqKCI5wGShcNNHVBzGSEga3E0CgYEA8lSLFRpqRbLbdRan
nWUO26OlBtJwadNCyotaB+dlWqqoNW1gXdvrLqIPash0+VuZyOc5i2TjL7WuHcP+
daITS0fR9fjCKJCLJd9mQRAMouInOgWn+3XQ3FWwnac+Oms1UkLQ3CD4T3mMCZgh
NiUfLd6wcdmIx6Di3SQhNcgEx6cCgYEA4ILSwvoZNMPsf3Pu8BQq5JgzjoWOkZcR
/VQla3+Zmp9piILfcmu5ZtLPn/hEa8c/Gz0TOkAOjrF93ZO/wwUpm0J6OeXtepgm
Qj0aZbYpBeVXXPJTjEYzApLZ2kRVS0Jgszao1ozaGS00KWPhleUxuW3F1t6ZtgUt
z0y8j7paA9MCgYAnKF4pFNGjXZl0vCptlozmFPVvusWLdCXQ0N/FczC/i+AOchQm
VokUSf+vw6CTQqgA/MwhqDXF2de+25Lcs0nm2j4lPwMwxtqVThoZ3VwNXfn0uEFC
svEDCZ22e/XkPrqpBj71syYswXlNe5WySCFXqipc20fs6iR+k34CUMXk9QKBgQCt
od/YdU4bNc9o/sNzH1XQ9zkgZ4BMbta14mcSTVwBbnpS3kxrlI6NzEpRANQButW+
fiyppzLa/sBGJmdvL0XvUqlur7lLT/J+1fzdXtU5PxaixrMA0KXQvrwHh0Oj8dER
qRtI2yQtxn0J6bpkkB41t4UDWaLtV/IG2eUXD1tgiwKBgQDYGqGQ+a3lgZXxTDpK
ETdzS5xBTq5duhDVjAXbZ1nbmZAw/781CYhS2ViWr2go/MnIZ0o5EZyfBo2ZxT2V
E6dPC2Uf48xSkru9oKt2oQJsuwKf0VObaQr1xM87hhgQF/nhIitjRisnzCdTWKXs
Ot3ZLcRlZaxBxPrmYxuPCoCO3A==
-----END PRIVATE KEY-----`

func main() {
	opts, err := parseArgs()
	if err != nil {
		fatal(err.Error())
	}
	// fmt.Println(opts.host)
	// fmt.Println(opts.port)
	if opts.version {
		fmt.Println(version)
		return
	}

	logger := log.NewFilterLogger(log.NewStdLogger(), opts.logLevel)

	// read configuration file
	var config *ClientConfig

	if opts.command != "qstart" {
		config, err = loadClientConfigFromFile(opts.config)
		if err != nil {
			fatal("configuration error: %s", err)
		}
	} else {

		var urlProto = "http"
		if opts.protocol == "tcp" {
			urlProto = "tcp"
		}
		config = &ClientConfig{
			ServerAddr: "tunnel.arumiot.com:5223",
			Tunnels: map[string]*Tunnel{
				"webui": &Tunnel{
					Protocol: opts.protocol,
					Host:     opts.host,
					Addr:     urlProto + "://localhost:" + strconv.Itoa(opts.port),
				},
			},
		}
	}

	switch opts.command {
	case "id":
		cert, err := tls.X509KeyPair([]byte(crt), []byte(key))
		// cert, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			fatal("failed to load key pair: %s", err)
		}
		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			fatal("failed to parse certificate: %s", err)
		}
		fmt.Println(id.New(x509Cert.Raw))

		return
	case "list":
		var names []string
		for n := range config.Tunnels {
			names = append(names, n)
		}

		sort.Strings(names)

		for _, n := range names {
			fmt.Println(n)
		}

		return
	case "start":
		tunnels := make(map[string]*Tunnel)
		for _, arg := range opts.args {
			t, ok := config.Tunnels[arg]
			if !ok {
				fatal("no such tunnel %q", arg)
			}
			tunnels[arg] = t
		}
		config.Tunnels = tunnels
	}

	if len(config.Tunnels) == 0 {
		fatal("no tunnels")
	}

	tlsconf, err := tlsConfig(config)
	if err != nil {
		fatal("failed to configure tls: %s", err)
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		fatal("failed to dump config: %s", err)
	}
	logger.Log("config", string(b))

	client, err := tunnel.NewClient(&tunnel.ClientConfig{
		ServerAddr:      config.ServerAddr,
		TLSClientConfig: tlsconf,
		Backoff:         expBackoff(config.Backoff),
		Tunnels:         tunnels(config.Tunnels),
		Proxy:           proxy(config.Tunnels, logger),
		Logger:          logger,
	})
	if err != nil {
		fatal("failed to create client: %s", err)
	}

	if err := client.Start(); err != nil {
		fatal("failed to start tunnels: %s", err)
	}
}

func tlsConfig(config *ClientConfig) (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(crt), []byte(key))
	// cert, err := tls.LoadX509KeyPair(config.TLSCrt, config.TLSKey)
	if err != nil {
		return nil, err
	}

	var roots *x509.CertPool
	if config.RootCA != "" {
		roots = x509.NewCertPool()
		rootPEM, err := ioutil.ReadFile(config.RootCA)
		if err != nil {
			return nil, err
		}
		if ok := roots.AppendCertsFromPEM(rootPEM); !ok {
			return nil, err
		}
	}

	host, _, err := net.SplitHostPort(config.ServerAddr)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ServerName:         host,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: roots == nil,
		RootCAs:            roots,
	}, nil
}

func expBackoff(c BackoffConfig) *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = c.Interval
	b.Multiplier = c.Multiplier
	b.MaxInterval = c.MaxInterval
	b.MaxElapsedTime = c.MaxTime

	return b
}

func tunnels(m map[string]*Tunnel) map[string]*proto.Tunnel {
	p := make(map[string]*proto.Tunnel)

	for name, t := range m {
		p[name] = &proto.Tunnel{
			Protocol: t.Protocol,
			Host:     t.Host,
			Auth:     t.Auth,
			Addr:     t.RemoteAddr,
		}
	}

	return p
}

func proxy(m map[string]*Tunnel, logger log.Logger) tunnel.ProxyFunc {
	httpURL := make(map[string]*url.URL)
	tcpAddr := make(map[string]string)

	for _, t := range m {
		switch t.Protocol {
		case proto.HTTP:
			u, err := url.Parse(t.Addr)
			if err != nil {
				fatal("invalid tunnel address: %s", err)
			}
			httpURL[t.Host] = u
		case proto.TCP, proto.TCP4, proto.TCP6:
			tcpAddr[t.RemoteAddr] = t.Addr
		case proto.SNI:
			tcpAddr[t.Host] = t.Addr
		}
	}

	return tunnel.Proxy(tunnel.ProxyFuncs{
		HTTP: tunnel.NewMultiHTTPProxy(httpURL, log.NewContext(logger).WithPrefix("proxy", "HTTP")).Proxy,
		TCP:  tunnel.NewMultiTCPProxy(tcpAddr, log.NewContext(logger).WithPrefix("proxy", "TCP")).Proxy,
	})
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}
