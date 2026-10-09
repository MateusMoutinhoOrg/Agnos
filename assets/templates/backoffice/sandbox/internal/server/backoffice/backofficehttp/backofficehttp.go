package backofficehttp

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// ContentSecurityPolicy is what a backoffice page may load: scripts only from
// the server itself (/backoffice/backoffice.js, never inline), the styles the
// templates carry inline, and no frame around it.
const ContentSecurityPolicy = "default-src 'none'; script-src 'self'; style-src 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; form-action 'self'; frame-ancestors 'none'; base-uri 'none'"

// HstsMaxAge is how long a browser keeps to https once it saw the backoffice
// over it, in seconds: a year.
const HstsMaxAge = 365 * 24 * 60 * 60

// IsIp tells whether ip is an IPv4 or an IPv6 address, as ParseIp reads it.
func IsIp(sandbox *api.Sandbox, ip string) (bool, error) {
	_, ok := ParseIp(sandbox, ip)
	return ok, nil
}

// ParseIp reads ip, an IPv4 address in dotted decimal (no leading zero in a
// byte) or an IPv6 one in any of its spellings — upper or lower case, `::`,
// leading zeros, a dotted IPv4 tail — into its 16 bytes, an IPv4 one as the
// IPv4-mapped IPv6 address (::ffff:a.b.c.d) so the two spellings of one
// address compare equal. ok is false for anything else, a zone included.
func ParseIp(sandbox *api.Sandbox, ip string) (address [16]byte, ok bool) {
	if v4, ok := parseIpv4(ip); ok {
		address[10], address[11] = 0xff, 0xff
		copy(address[12:], v4[:])
		return address, true
	}
	return parseIpv6(ip)
}

// parseIpv4 reads a dotted decimal IPv4 address.
func parseIpv4(text string) (address [4]byte, ok bool) {
	part, value, digits := 0, 0, 0
	for i := 0; i <= len(text); i++ {
		if i == len(text) || text[i] == '.' {
			if digits == 0 || part > 3 {
				return address, false
			}
			address[part] = byte(value)
			part, value, digits = part+1, 0, 0
			continue
		}
		char := text[i]
		if char < '0' || char > '9' || (digits == 1 && value == 0) {
			return address, false
		}
		value = value*10 + int(char-'0')
		digits++
		if value > 255 {
			return address, false
		}
	}
	return address, part == 4
}

// parseIpv6 reads an IPv6 address: eight groups of up to four hex digits, a
// run of them written `::` once at most, the last two optionally a dotted
// IPv4 address.
func parseIpv6(text string) (address [16]byte, ok bool) {
	groups := [8]uint16{}
	count, gap := 0, -1
	i := 0
	if len(text) >= 2 && text[0] == ':' && text[1] == ':' {
		gap, i = 0, 2
	} else if len(text) > 0 && text[0] == ':' {
		return address, false
	}
	for i < len(text) {
		if count == 8 {
			return address, false
		}
		end := i
		for end < len(text) && text[end] != ':' {
			end++
		}
		field := text[i:end]
		if end == len(text) && count <= 6 && containsDot(field) {
			v4, ok := parseIpv4(field)
			if !ok {
				return address, false
			}
			groups[count] = uint16(v4[0])<<8 | uint16(v4[1])
			groups[count+1] = uint16(v4[2])<<8 | uint16(v4[3])
			count += 2
			i = end
			break
		}
		value, ok := parseHexGroup(field)
		if !ok {
			return address, false
		}
		groups[count] = value
		count++
		i = end
		if i == len(text) {
			break
		}
		i++
		if i < len(text) && text[i] == ':' {
			if gap >= 0 {
				return address, false
			}
			gap = count
			i++
		} else if i == len(text) {
			return address, false
		}
	}
	if gap < 0 && count != 8 || gap >= 0 && count > 7 {
		return address, false
	}
	if gap >= 0 {
		shift := 8 - count
		for index := count - 1; index >= gap; index-- {
			groups[index+shift] = groups[index]
			groups[index] = 0
		}
	}
	for index, group := range groups {
		address[2*index], address[2*index+1] = byte(group>>8), byte(group)
	}
	return address, true
}

// parseHexGroup reads one to four hex digits.
func parseHexGroup(field string) (uint16, bool) {
	if len(field) == 0 || len(field) > 4 {
		return 0, false
	}
	value := uint16(0)
	for i := 0; i < len(field); i++ {
		char := field[i]
		switch {
		case char >= '0' && char <= '9':
			value = value<<4 | uint16(char-'0')
		case char >= 'a' && char <= 'f':
			value = value<<4 | uint16(char-'a'+10)
		case char >= 'A' && char <= 'F':
			value = value<<4 | uint16(char-'A'+10)
		default:
			return 0, false
		}
	}
	return value, true
}

// containsDot tells whether text holds a '.'.
func containsDot(text string) bool {
	for i := 0; i < len(text); i++ {
		if text[i] == '.' {
			return true
		}
	}
	return false
}

// isMapped tells whether address is an IPv4-mapped one, ::ffff:a.b.c.d.
func isMapped(address [16]byte) bool {
	for i := 0; i < 10; i++ {
		if address[i] != 0 {
			return false
		}
	}
	return address[10] == 0xff && address[11] == 0xff
}

// FormatIp spells address the one way ParseIp's caller reads back: dotted
// decimal for an IPv4 one, the RFC 5952 form for an IPv6 one — lower case, no
// leading zero, the longest run of zero groups written `::`.
func FormatIp(sandbox *api.Sandbox, address [16]byte) string {
	strings := sandbox.Deps.StringsDeps
	if isMapped(address) {
		return strings.FormatInt(int64(address[12]), 10) + "." + strings.FormatInt(int64(address[13]), 10) + "." +
			strings.FormatInt(int64(address[14]), 10) + "." + strings.FormatInt(int64(address[15]), 10)
	}
	groups := [8]int64{}
	for index := range groups {
		groups[index] = int64(address[2*index])<<8 | int64(address[2*index+1])
	}
	bestStart, bestLength := -1, 1
	for start := 0; start < 8; {
		if groups[start] != 0 {
			start++
			continue
		}
		end := start
		for end < 8 && groups[end] == 0 {
			end++
		}
		if end-start > bestLength {
			bestStart, bestLength = start, end-start
		}
		start = end
	}
	text := ""
	for index := 0; index < 8; index++ {
		if index == bestStart {
			text += "::"
			index += bestLength - 1
			continue
		}
		if text != "" && !strings.HasSuffix(text, ":") {
			text += ":"
		}
		text += strings.FormatInt(groups[index], 16)
	}
	return text
}

// CanonicalIpOrRange is entry — an ip, or a CIDR range "address/bits" — in
// the one spelling FormatIp gives it, a range's address cut to its prefix;
// ok is false when entry is neither.
func CanonicalIpOrRange(sandbox *api.Sandbox, entry string) (string, bool) {
	address, bits, ranged, ok := parseIpOrRange(sandbox, entry)
	if !ok {
		return "", false
	}
	if !ranged {
		return FormatIp(sandbox, address), true
	}
	shown := bits
	if isMapped(address) {
		shown = bits - 96
	}
	return FormatIp(sandbox, masked(address, bits)) + "/" + sandbox.Deps.StringsDeps.FormatInt(int64(shown), 10), true
}

// IpMatches tells whether ip is entry — an address, compared as parsed — or
// falls within it, a CIDR range.
func IpMatches(sandbox *api.Sandbox, entry string, ip string) bool {
	target, ok := ParseIp(sandbox, ip)
	if !ok {
		return false
	}
	address, bits, _, ok := parseIpOrRange(sandbox, entry)
	if !ok {
		return false
	}
	return masked(address, bits) == masked(target, bits)
}

// parseIpOrRange reads entry as an ip — bits 128 — or a CIDR range, whose
// bits count over the 16 bytes ParseIp answers: an IPv4 range's own plus 96.
func parseIpOrRange(sandbox *api.Sandbox, entry string) (address [16]byte, bits int, ranged bool, ok bool) {
	strings := sandbox.Deps.StringsDeps
	cut := strings.LastIndex(entry, "/")
	if cut < 0 {
		address, ok = ParseIp(sandbox, entry)
		return address, 128, false, ok
	}
	address, ok = ParseIp(sandbox, entry[:cut])
	if !ok {
		return address, 0, false, false
	}
	size := entry[cut+1:]
	if size == "" || len(size) > 3 || (len(size) > 1 && size[0] == '0') {
		return address, 0, false, false
	}
	parsed, err := strings.ParseInt(size, 10, 64)
	limit := int64(128)
	if _, v4 := parseIpv4(entry[:cut]); v4 {
		limit = 32
	}
	if err != nil || parsed < 0 || parsed > limit {
		return address, 0, false, false
	}
	bits = int(parsed)
	if limit == 32 {
		bits += 96
	}
	return address, bits, true, true
}

// masked is address with every bit past the first bits cleared.
func masked(address [16]byte, bits int) [16]byte {
	for index := range address {
		keep := bits - 8*index
		switch {
		case keep >= 8:
		case keep <= 0:
			address[index] = 0
		default:
			address[index] &= byte(0xff << (8 - keep))
		}
	}
	return address
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
	strings := sandbox.Deps.StringsDeps
	entries := strings.Split(forwardedFor, ",")
	last, valid := ParseIp(sandbox, strings.TrimSpace(entries[len(entries)-1]))
	if !valid {
		return peer
	}
	return FormatIp(sandbox, last)
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
		response.SetHeader("Strict-Transport-Security", sandbox.Deps.StdDeps.Sprintf("max-age=%d", HstsMaxAge))
	}
}

// BaseSecurityHeaders sets on response the headers every answer outside the
// backoffice carries — the application's pages and files — unless the route
// answering sets its own: no content-type sniffing, no framing by another
// site, and no full url sent as referrer to another site. It sets no
// Content-Security-Policy: what a page of the application may load is the
// application's to say.
func BaseSecurityHeaders(sandbox *api.Sandbox, response *serverdeps.Response) {
	response.SetHeader("X-Content-Type-Options", "nosniff")
	response.SetHeader("X-Frame-Options", "SAMEORIGIN")
	response.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin")
}

// IsBackofficePath tells whether path is one of the backoffice's own: /admin
// or /api/admin, or anything under either.
func IsBackofficePath(sandbox *api.Sandbox, path string) bool {
	strings := sandbox.Deps.StringsDeps
	for _, root := range []string{"/admin", "/api/admin"} {
		if path == root || strings.HasPrefix(path, root+"/") {
			return true
		}
	}
	return false
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
	strings := sandbox.Deps.StringsDeps
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
	strings := sandbox.Deps.StringsDeps
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
