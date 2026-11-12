// Command certdays prints days until TLS certificate expiry for each host and
// exits non-zero when any certificate is close to expiry or cannot be checked.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PrinceAdvant/certdays"
)

func main() {
	warn := flag.Int("warn", 14, "exit 1 if any certificate expires within this many days")
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: certdays [-warn 14] host[:port] [host...]")
		os.Exit(2)
	}
	bad := false
	for _, h := range flag.Args() {
		r := certdays.Check(h, *timeout)
		if r.Err != nil {
			fmt.Printf("%s error: %v\n", r.Host, r.Err)
			bad = true
			continue
		}
		status := "ok"
		if r.DaysLeft < *warn {
			status = "RENEW"
			bad = true
		}
		fmt.Printf("%s %s %d days (%s) issuer=%q names=%s\n", r.Host, status, r.DaysLeft, r.NotAfter.Format("2006-01-02"), r.Issuer, strings.Join(r.Names, ","))
	}
	if bad {
		os.Exit(1)
	}
}
