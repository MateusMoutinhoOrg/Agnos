package stringsdeps

import (
	"regexp"
	"strconv"
	"strings"

	stringsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/stringsdeps"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// Bind fills deps.Deps.Stringsdeps with the standard library's strings and
// strconv. Every field is a straight delegation: the contract restates the
// standard library api so the sandbox can call it without importing it.
func Bind(deps *deps.Deps) {
	deps.Stringsdeps = stringsdeps.Lib{
		TrimSpace: func(s string) string {
			return strings.TrimSpace(s)
		},
		Trim: func(s string, cutset string) string {
			return strings.Trim(s, cutset)
		},
		TrimLeft: func(s string, cutset string) string {
			return strings.TrimLeft(s, cutset)
		},
		TrimRight: func(s string, cutset string) string {
			return strings.TrimRight(s, cutset)
		},
		TrimPrefix: func(s string, prefix string) string {
			return strings.TrimPrefix(s, prefix)
		},
		TrimSuffix: func(s string, suffix string) string {
			return strings.TrimSuffix(s, suffix)
		},
		HasPrefix: func(s string, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		HasSuffix: func(s string, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
		Contains: func(s string, substr string) bool {
			return strings.Contains(s, substr)
		},
		ContainsAny: func(s string, chars string) bool {
			return strings.ContainsAny(s, chars)
		},
		LastIndex: func(s string, substr string) int {
			return strings.LastIndex(s, substr)
		},
		Count: func(s string, substr string) int {
			return strings.Count(s, substr)
		},
		Split: func(s string, sep string) []string {
			return strings.Split(s, sep)
		},
		Join: func(elems []string, sep string) string {
			return strings.Join(elems, sep)
		},
		Fields: func(s string) []string {
			return strings.Fields(s)
		},
		FieldsFunc: func(s string, f func(rune) bool) []string {
			return strings.FieldsFunc(s, f)
		},
		Repeat: func(s string, count int) string {
			return strings.Repeat(s, count)
		},
		ReplaceAll: func(s string, old string, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		ToUpper: func(s string) string {
			return strings.ToUpper(s)
		},
		ToLower: func(s string) string {
			return strings.ToLower(s)
		},
		Quote: func(s string) string {
			return strconv.Quote(s)
		},
		MatchPattern: func(pattern string, s string) (bool, error) {
			return regexp.MatchString(pattern, s)
		},
		Atoi: func(s string) (int, error) {
			return strconv.Atoi(s)
		},
		ParseInt: func(s string, base int, bit_size int) (int64, error) {
			return strconv.ParseInt(s, base, bit_size)
		},
		ParseFloat: func(s string, bit_size int) (float64, error) {
			return strconv.ParseFloat(s, bit_size)
		},
		FormatInt: func(value int64, base int) string {
			return strconv.FormatInt(value, base)
		},
		FormatFloat: func(value float64, format byte, precision int, bit_size int) string {
			return strconv.FormatFloat(value, format, precision, bit_size)
		},
	}
}
