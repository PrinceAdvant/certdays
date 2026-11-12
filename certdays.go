// Package certdays reports how many days remain before the TLS certificate
// served by a host expires, and who issued it.
package certdays

import (
	"crypto/tls"
	"net"
	"strings"
	"time"
)

// Result describes the leaf certificate presented by a host.
type Result struct {
	Host     string
	Issuer   string
	Names    []string
	NotAfter time.Time
	DaysLeft int
	Err      error
}

// DaysUntil returns whole days between now and t, rounded down.
func DaysUntil(now, t time.Time) int {
	return int(t.Sub(now).Hours() / 24)
}

// Normalize turns "example.com" into "example.com:443" and keeps explicit ports.
func Normalize(target string) (addr, serverName string) {
	target = strings.TrimPrefix(strings.TrimPrefix(target, "https://"), "http://")
	target = strings.SplitN(target, "/", 2)[0]
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		host, port = target, "443"
	}
	return net.JoinHostPort(host, port), host
}

// Check connects to target and inspects the certificate chain. Verification
// is performed so that a wrong name or untrusted issuer is reported as an error.
func Check(target string, timeout time.Duration) Result {
	addr, name := Normalize(target)
	r := Result{Host: addr}
	d := &net.Dialer{Timeout: timeout}
	conn, err := tls.DialWithDialer(d, "tcp", addr, &tls.Config{ServerName: name})
	if err != nil {
		r.Err = err
		return r
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return r
	}
	leaf := certs[0]
	r.Issuer = leaf.Issuer.CommonName
	if len(leaf.Issuer.Organization) > 0 {
		r.Issuer = leaf.Issuer.Organization[0] + " " + r.Issuer
	}
	r.Names = leaf.DNSNames
	r.NotAfter = leaf.NotAfter
	r.DaysLeft = DaysUntil(time.Now(), leaf.NotAfter)
	return r
}
