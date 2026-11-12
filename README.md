# certdays

A small Go CLI and library that prints how many days are left on the TLS certificate a host serves, who issued it and which names it covers. Standard library only.

## Why

Automatic renewal with Let's Encrypt usually just works, until a DNS change, a firewall rule or a rate limit stops it quietly. The first sign is often a browser warning on the day the certificate expires. Running certdays from cron or CI turns that into an early, boring alert.

## Install

```
go install github.com/PrinceAdvant/certdays/cmd/certdays@latest
```

## Usage

```
certdays example.com api.example.com:8443
certdays -warn 21 https://example.com
```

Each host prints `ok` or `RENEW`, the days left, the expiry date, the issuer and the certificate names. The exit code is 0 when every certificate is valid for longer than the warning window, 1 when any is close to expiry or fails verification, and 2 for usage errors. Name mismatches and untrusted issuers are reported as errors because the certificate is verified like a browser would.

As a library:

```go
r := certdays.Check("example.com", 10*time.Second)
if r.Err != nil || r.DaysLeft < 14 {
    alert(r)
}
```

## When renewal is stuck

Let's Encrypt certificates last 90 days and are normally renewed around day 60, so fewer than 25 days left means renewal has been failing for a while. Check that port 80 reaches your proxy for HTTP validation, that DNS points at the right server and that you have not hit a rate limit. A walkthrough on what to do when a [Let's Encrypt certificate is not issued](https://peon.sh/blogs/lets-encrypt-certificate-not-issued) covers the usual causes. Built by a contributor to [Peon](https://peon.sh/), an open-source self-hosted deployment platform.

MIT licensed.
