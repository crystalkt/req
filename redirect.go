package req

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type RedirectPolicy func(req *http.Request, via []*http.Request) error

// MaxRedirectPolicy specifies the max number of redirect
func MaxRedirectPolicy(noOfRedirect int) RedirectPolicy {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= noOfRedirect {
			return fmt.Errorf("stopped after %d redirects", noOfRedirect)
		}
		return nil
	}
}

// NoRedirectPolicy disable redirect behaviour
func NoRedirectPolicy() RedirectPolicy {
	return func(req *http.Request, via []*http.Request) error {
		return errors.New("auto redirect is disabled")
	}
}

func SameDomainRedirectPolicy() RedirectPolicy {
	return func(req *http.Request, via []*http.Request) error {
		if getDomain(req.URL.Host) != getDomain(via[0].URL.Host) {
			return errors.New("different domain name is not allowed")
		}
		return nil
	}
}

// SameHostRedirectPolicy allows redirect only if the redirected host
// is the same as original host, e.g. redirect to "www.imroc.cc" from
// "imroc.cc" is not the allowed.
func SameHostRedirectPolicy() RedirectPolicy {
	return func(req *http.Request, via []*http.Request) error {
		if getHostname(req.URL.Host) != getHostname(via[0].URL.Host) {
			return errors.New("different host name is not allowed")
		}
		return nil
	}
}

// AllowedHostRedirectPolicy allows redirect only if the redirected host
// match one of the host that specified.
func AllowedHostRedirectPolicy(hosts ...string) RedirectPolicy {
	m := make(map[string]bool)
	for _, h := range hosts {
		m[strings.ToLower(getHostname(h))] = true
	}

	return func(req *http.Request, via []*http.Request) error {
		if _, ok := m[getHostname(req.URL.Host)]; !ok {
			return errors.New("redirect host is not allowed")
		}
		return nil
	}
}

// AllowedDomainRedirectPolicy allows redirect only if the redirected domain
// match one of the domain that specified.
func AllowedDomainRedirectPolicy(hosts ...string) RedirectPolicy {
	domains := make(map[string]bool)
	for _, h := range hosts {
		domains[strings.ToLower(getDomain(h))] = true
	}

	return func(req *http.Request, via []*http.Request) error {
		if _, ok := domains[getDomain(req.URL.Host)]; !ok {
			return errors.New("redirect domain is not allowed")
		}
		return nil
	}
}

func getHostname(host string) (hostname string) {
	if strings.Index(host, ":") > 0 {
		host, _, _ = net.SplitHostPort(host)
	}
	hostname = strings.ToLower(host)
	return
}

func getDomain(host string) string {
	host = getHostname(host)
	ss := strings.Split(host, ".")
	if len(ss) < 3 {
		return host
	}
	ss = ss[1:]
	return strings.Join(ss, ".")
}
