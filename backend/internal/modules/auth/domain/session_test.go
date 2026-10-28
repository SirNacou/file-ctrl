package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"testing"
	"time"
)

var lowerHex = regexp.MustCompile(`^[0-9a-f]+$`)

func mustNewSessionToken(t *testing.T) SessionToken {
	t.Helper()

	tok, err := NewSessionToken()
	if err != nil {
		t.Fatalf("NewSessionToken() unexpected error: %v", err)
	}
	return tok
}

// ---------------------------------------------------------------------------
// NewSessionToken – raw token
// ---------------------------------------------------------------------------

func TestNewSessionToken_RawFormat(t *testing.T) {
	tok := mustNewSessionToken(t)

	if want := TokenLength * 2; len(tok.Raw) != want {
		t.Errorf("len(Raw) = %d, want %d", len(tok.Raw), want)
	}

	if !lowerHex.MatchString(tok.Raw) {
		t.Errorf("Raw must be lowercase hex, got %q", tok.Raw)
	}

	decoded, err := hex.DecodeString(tok.Raw)
	if err != nil {
		t.Fatalf("Raw is not valid hex: %v", err)
	}
	if len(decoded) != TokenLength {
		t.Errorf("decoded Raw length = %d bytes, want %d", len(decoded), TokenLength)
	}
}

func TestNewSessionToken_RawIsNotAllZeros(t *testing.T) {
	tok := mustNewSessionToken(t)

	decoded, err := hex.DecodeString(tok.Raw)
	if err != nil {
		t.Fatalf("Raw is not valid hex: %v", err)
	}

	for _, b := range decoded {
		if b != 0 {
			return
		}
	}
	t.Error("Raw decoded to all zero bytes; random source looks broken")
}

func TestNewSessionToken_UniqueAcrossCalls(t *testing.T) {
	const n = 1000

	rawSeen := make(map[string]struct{}, n)
	hashSeen := make(map[string]struct{}, n)

	for i := range n {
		tok := mustNewSessionToken(t)

		if _, dup := rawSeen[tok.Raw]; dup {
			t.Fatalf("duplicate Raw token after %d iterations: %s", i, tok.Raw)
		}
		if _, dup := hashSeen[tok.Hash]; dup {
			t.Fatalf("duplicate Hash after %d iterations: %s", i, tok.Hash)
		}

		rawSeen[tok.Raw] = struct{}{}
		hashSeen[tok.Hash] = struct{}{}
	}
}

// ---------------------------------------------------------------------------
// NewSessionToken – hash
// ---------------------------------------------------------------------------

func TestNewSessionToken_HashFormat(t *testing.T) {
	tok := mustNewSessionToken(t)

	// SHA-256 -> 32 bytes -> 64 hex chars.
	if len(tok.Hash) != sha256.Size*2 {
		t.Errorf("len(Hash) = %d, want %d", len(tok.Hash), sha256.Size*2)
	}
	if !lowerHex.MatchString(tok.Hash) {
		t.Errorf("Hash must be lowercase hex, got %q", tok.Hash)
	}
}

func TestNewSessionToken_HashMatchesRaw(t *testing.T) {
	tok := mustNewSessionToken(t)

	// Compute independently of the production helper.
	sum := sha256.Sum256([]byte(tok.Raw))
	want := hex.EncodeToString(sum[:])

	if tok.Hash != want {
		t.Errorf("Hash = %q, want sha256(Raw) = %q", tok.Hash, want)
	}

	if tok.Hash != hashSessionToken(tok.Raw) {
		t.Error("Hash must equal hashSessionToken(Raw)")
	}
}

func TestNewSessionToken_HashDiffersFromRaw(t *testing.T) {
	tok := mustNewSessionToken(t)

	if tok.Hash == tok.Raw {
		t.Error("Hash must not equal Raw; the raw token must not be stored as-is")
	}
}

// ---------------------------------------------------------------------------
// NewSessionToken – expiry
// ---------------------------------------------------------------------------

func TestNewSessionToken_ExpiresAfterTTL(t *testing.T) {
	before := time.Now().UTC()
	tok := mustNewSessionToken(t)
	after := time.Now().UTC()

	earliest := before.Add(SessionTTL)
	latest := after.Add(SessionTTL)

	if tok.ExpiresAt.Before(earliest) || tok.ExpiresAt.After(latest) {
		t.Errorf("ExpiresAt = %v, want within [%v, %v]", tok.ExpiresAt, earliest, latest)
	}
}

func TestNewSessionToken_ExpiresAtIsUTC(t *testing.T) {
	tok := mustNewSessionToken(t)

	if loc := tok.ExpiresAt.Location(); loc != time.UTC {
		t.Errorf("ExpiresAt location = %v, want UTC", loc)
	}
}

func TestNewSessionToken_ExpiresInFuture(t *testing.T) {
	tok := mustNewSessionToken(t)

	if !tok.ExpiresAt.After(time.Now()) {
		t.Errorf("ExpiresAt = %v should be in the future", tok.ExpiresAt)
	}
}

func TestSessionTTL_IsSevenDays(t *testing.T) {
	if want := 7 * 24 * time.Hour; SessionTTL != want {
		t.Errorf("SessionTTL = %v, want %v", SessionTTL, want)
	}
}

func TestTokenLength_Is32Bytes(t *testing.T) {
	// 32 bytes = 256 bits of entropy. Guard against accidental weakening.
	if TokenLength < 32 {
		t.Errorf("TokenLength = %d, want at least 32 bytes", TokenLength)
	}
}

// ---------------------------------------------------------------------------
// hashSessionToken
// ---------------------------------------------------------------------------

func TestHashSessionToken_KnownVectors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:  "abc",
			input: "abc",
			want:  "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
		{
			name:  "hello",
			input: "hello",
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashSessionToken(tt.input); got != tt.want {
				t.Errorf("hashSessionToken(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHashSessionToken_Deterministic(t *testing.T) {
	const raw = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	first := hashSessionToken(raw)
	for i := range 10 {
		if got := hashSessionToken(raw); got != first {
			t.Fatalf("call %d returned %q, want %q", i, got, first)
		}
	}
}

func TestHashSessionToken_DifferentInputsDifferentHashes(t *testing.T) {
	inputs := []string{
		"token-a",
		"token-b",
		"Token-a",  // case difference
		"token-a ", // trailing space
		" token-a", // leading space
		"",
	}

	seen := make(map[string]string, len(inputs))
	for _, in := range inputs {
		h := hashSessionToken(in)
		if prev, dup := seen[h]; dup {
			t.Errorf("inputs %q and %q produced the same hash %q", prev, in, h)
		}
		seen[h] = in
	}
}

func TestHashSessionToken_OutputFormat(t *testing.T) {
	h := hashSessionToken("anything")

	if len(h) != 64 {
		t.Errorf("len = %d, want 64", len(h))
	}
	if !lowerHex.MatchString(h) {
		t.Errorf("output must be lowercase hex, got %q", h)
	}
}

// ---------------------------------------------------------------------------
// Struct behaviour
// ---------------------------------------------------------------------------

// A client presenting the raw token later must be able to recompute the
// stored hash and find the matching session.
func TestSessionToken_LookupRoundTrip(t *testing.T) {
	tok := mustNewSessionToken(t)

	presented := tok.Raw // what the client sends back
	if got := hashSessionToken(presented); got != tok.Hash {
		t.Errorf("hash of presented token = %q, want stored hash %q", got, tok.Hash)
	}

	// A tampered token must not match.
	tampered := presented[:len(presented)-1] + flipHexChar(presented[len(presented)-1])
	if hashSessionToken(tampered) == tok.Hash {
		t.Error("tampered token must not hash to the stored hash")
	}
}

func flipHexChar(c byte) string {
	if c == '0' {
		return "1"
	}
	return "0"
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkNewSessionToken(b *testing.B) {
	for b.Loop() {
		if _, err := NewSessionToken(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashSessionToken(b *testing.B) {
	const raw = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	for b.Loop() {
		_ = hashSessionToken(raw)
	}
}
