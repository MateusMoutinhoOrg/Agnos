package backofficeguard

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// ContentSecurityPolicy is what a backoffice page may load: scripts only from
// the server itself (/admin/backoffice.js, never inline), the styles the
// templates carry inline, and no frame around it.
const ContentSecurityPolicy = "default-src 'none'; script-src 'self'; style-src 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; form-action 'self'; frame-ancestors 'none'; base-uri 'none'"

// HstsMaxAge is how long a browser keeps to https once it saw the backoffice
// over it, in seconds: a year.
const HstsMaxAge = 365 * 24 * 60 * 60

// ipv4Pattern is a dotted IPv4 address, every byte within 0-255.
const ipv4Pattern = `^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`

// ipv6Pattern is an IPv6 address in lower case, groups of up to four hex
// digits around colons, a `::` included.
const ipv6Pattern = `^[0-9a-f]{0,4}(:[0-9a-f]{0,4}){2,7}$`

// IsIp tells whether ip, already lower case, is an IPv4 or an IPv6 address.
func IsIp(sandbox *api.Sandbox, ip string) (bool, error) {
	ipv4, err := sandbox.Deps.Stringsdeps.MatchPattern(ipv4Pattern, ip)
	if err != nil || ipv4 {
		return ipv4, err
	}
	return sandbox.Deps.Stringsdeps.MatchPattern(ipv6Pattern, ip)
}

// ClientIp is the ip a request came from, given peer, the ip of its
// connection, and forwardedFor, every X-Forwarded-For it carried joined by
// commas. Unless start-server trusts X-Forwarded-For, it is peer: anyone may
// send that header. When it does, the server sits behind one reverse proxy
// that appends the ip it was reached from, so the last entry is the client's
// and every entry before it is whatever the client sent; a last entry that is
// not an ip, or none at all, falls back to peer.
func ClientIp(sandbox *api.Sandbox, peer string, forwardedFor string) string {
	if !sandbox.Config.AllowXForwardedFor {
		return peer
	}
	strings := sandbox.Deps.Stringsdeps
	entries := strings.Split(forwardedFor, ",")
	last := strings.ToLower(strings.TrimSpace(entries[len(entries)-1]))
	valid, err := IsIp(sandbox, last)
	if err != nil || !valid {
		return peer
	}
	return last
}

// SecurityHeaders sets on response the headers every backoffice answer
// carries, page or JSON: what the page may load and who may frame it, no
// content-type sniffing, no referrer to another site, no caching, and — unless
// start-server serves plain http — that the browser keeps to https.
func SecurityHeaders(sandbox *api.Sandbox, response *serverdeps.Response) {
	response.SetHeader("Content-Security-Policy", ContentSecurityPolicy)
	response.SetHeader("X-Frame-Options", "DENY")
	response.SetHeader("X-Content-Type-Options", "nosniff")
	response.SetHeader("Referrer-Policy", "same-origin")
	response.SetHeader("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
	response.SetHeader("Cross-Origin-Opener-Policy", "same-origin")
	response.SetHeader("Cache-Control", "no-store")
	if !sandbox.Config.InsecureHttp {
		response.SetHeader("Strict-Transport-Security", sandbox.Deps.Std.Sprintf("max-age=%d", HstsMaxAge))
	}
}

// SameOrigin tells whether a request whose Origin header is origin was sent
// by a page of host, the host it was sent to. A browser sends Origin on every
// form post, so one naming another host is a page of another site posting
// here; one spelled "null" comes from an opaque origin and is refused too. A
// request without Origin — a link followed, a client that is no browser — is
// let through: the session cookie, SameSite=Strict, is what keeps another
// site's navigation from carrying a session.
func SameOrigin(sandbox *api.Sandbox, origin string, host string) bool {
	if origin == "" {
		return true
	}
	strings := sandbox.Deps.Stringsdeps
	for _, scheme := range []string{"https://", "http://"} {
		if strings.HasPrefix(origin, scheme) {
			return host != "" && strings.ToLower(origin[len(scheme):]) == strings.ToLower(host)
		}
	}
	return false
}

// ListensEverywhere tells whether start-server's addr — a port, a range of
// them, either one behind a host — listens on every interface of the machine
// rather than on one address: no host, 0.0.0.0 or [::]. A server trusting
// X-Forwarded-For has to be reachable by its proxy alone, so it is warned.
func ListensEverywhere(sandbox *api.Sandbox, addr string) bool {
	strings := sandbox.Deps.Stringsdeps
	head := ""
	cut := strings.LastIndex(addr, ":")
	if cut >= 0 {
		head = addr[:cut]
	}
	host := head
	hostCut := strings.LastIndex(head, ":")
	if digits(sandbox, head[hostCut+1:]) {
		host = ""
		if hostCut >= 0 {
			host = head[:hostCut]
		}
	}
	return host == "" || host == "0.0.0.0" || host == "[::]" || host == "::"
}

// digits tells whether text is a non-empty run of decimal digits.
func digits(sandbox *api.Sandbox, text string) bool {
	if text == "" {
		return false
	}
	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return false
		}
	}
	return true
}
