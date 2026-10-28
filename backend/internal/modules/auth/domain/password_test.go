package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/argon2"
)

// buildHash produces an encoded argon2id hash with caller-chosen (cheap)
// parameters so verification tests stay fast.
func buildHash(
	t *testing.T,
	password string,
	memory, iterations uint32,
	parallelism uint8,
	salt []byte,
	keyLen uint32,
) string {
	t.Helper()

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

func mustPassword(t *testing.T, raw string) Password {
	t.Helper()

	p, err := NewPassword(raw)
	if err != nil {
		t.Fatalf("NewPassword(%q) unexpected error: %v", raw, err)
	}
	return p
}

// ---------------------------------------------------------------------------
// NewPassword
// ---------------------------------------------------------------------------

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  error
		wantSame bool // expect value to be stored as given
	}{
		{"empty", "", ErrPasswordTooShort, false},
		{
			"one below minimum",
			strings.Repeat("a", MinimumPasswordLength-1),
			ErrPasswordTooShort,
			false,
		},
		{"exactly minimum", strings.Repeat("a", MinimumPasswordLength), nil, true},
		{"typical", "correct-horse-battery-staple", nil, true},
		{"exactly maximum", strings.Repeat("a", MaximumPasswordLength), nil, true},
		{
			"one above maximum",
			strings.Repeat("a", MaximumPasswordLength+1),
			ErrPasswordTooLong,
			false,
		},
		{"far above maximum", strings.Repeat("a", 10_000), ErrPasswordTooLong, false},
		{"whitespace only but long enough", "        ", nil, true},
		{"unicode within byte limits", "pässwörd!", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPassword(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				if got != (Password{}) {
					t.Errorf("expected zero Password on error, got %+v", got)
				}
				return
			}

			if tt.wantSame && got.value != tt.input {
				t.Errorf("value = %q, want %q", got.value, tt.input)
			}
		})
	}
}

// Length is measured in bytes (len), not runes. This test documents that
// behavior: 4 two-byte runes == 8 bytes, which satisfies the minimum even
// though it is only 4 characters.
func TestNewPassword_LengthIsMeasuredInBytes(t *testing.T) {
	if _, err := NewPassword("éééé"); err != nil {
		t.Errorf("4 runes / 8 bytes should be accepted, got %v", err)
	}

	// 65 two-byte runes = 130 bytes > 128, though only 65 characters.
	if _, err := NewPassword(strings.Repeat("é", 65)); !errors.Is(err, ErrPasswordTooLong) {
		t.Errorf("130 bytes should be too long, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// HashPassword
// ---------------------------------------------------------------------------

func TestHashPassword_Format(t *testing.T) {
	p := mustPassword(t, "super-secret-pw")

	encoded, err := p.HashPassword()
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		t.Fatalf("expected 6 parts, got %d: %q", len(parts), encoded)
	}

	if parts[0] != "" {
		t.Errorf("expected leading empty segment, got %q", parts[0])
	}
	if parts[1] != "argon2id" {
		t.Errorf("algorithm = %q, want argon2id", parts[1])
	}
	if want := fmt.Sprintf("v=%d", argon2.Version); parts[2] != want {
		t.Errorf("version = %q, want %q", parts[2], want)
	}

	wantParams := fmt.Sprintf("m=%d,t=%d,p=%d",
		DefaultArgon2Params.Memory,
		DefaultArgon2Params.Iterations,
		DefaultArgon2Params.Parallelism,
	)
	if parts[3] != wantParams {
		t.Errorf("params = %q, want %q", parts[3], wantParams)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		t.Fatalf("salt is not valid raw base64: %v", err)
	}
	if uint32(len(salt)) != DefaultArgon2Params.SaltLength {
		t.Errorf("salt length = %d, want %d", len(salt), DefaultArgon2Params.SaltLength)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		t.Fatalf("hash is not valid raw base64: %v", err)
	}
	if uint32(len(hash)) != DefaultArgon2Params.KeyLength {
		t.Errorf("hash length = %d, want %d", len(hash), DefaultArgon2Params.KeyLength)
	}

	if strings.Contains(encoded, "=") && strings.Contains(parts[4]+parts[5], "=") {
		t.Error("base64 segments should be unpadded (RawStdEncoding)")
	}
}

func TestHashPassword_DoesNotLeakPlaintext(t *testing.T) {
	raw := "super-secret-pw"
	p := mustPassword(t, raw)

	encoded, err := p.HashPassword()
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	if strings.Contains(encoded, raw) {
		t.Error("encoded hash must not contain the plaintext password")
	}
}

func TestHashPassword_UniqueSaltPerCall(t *testing.T) {
	p := mustPassword(t, "same-password-twice")

	h1, err := p.HashPassword()
	if err != nil {
		t.Fatalf("first HashPassword() error: %v", err)
	}
	h2, err := p.HashPassword()
	if err != nil {
		t.Fatalf("second HashPassword() error: %v", err)
	}

	if h1 == h2 {
		t.Error("hashing the same password twice must produce different encodings (random salt)")
	}

	// Both must still verify.
	for i, h := range []string{h1, h2} {
		ok, err := p.VerifyPassword(h)
		if err != nil || !ok {
			t.Errorf("hash #%d failed to verify: ok=%v err=%v", i+1, ok, err)
		}
	}
}

func TestHashPassword_DoesNotMutatePassword(t *testing.T) {
	p := mustPassword(t, "immutable-check")
	before := p

	if _, err := p.HashPassword(); err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	if p != before {
		t.Errorf("Password mutated by HashPassword: before=%+v after=%+v", before, p)
	}
}

// ---------------------------------------------------------------------------
// VerifyPassword – happy path & mismatches
// ---------------------------------------------------------------------------

func TestVerifyPassword_RoundTrip(t *testing.T) {
	raw := "round-trip-password-123!"
	p := mustPassword(t, raw)

	encoded, err := p.HashPassword()
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	ok, err := p.VerifyPassword(encoded)
	if err != nil {
		t.Fatalf("VerifyPassword() error: %v", err)
	}
	if !ok {
		t.Error("expected correct password to verify")
	}

	// A separately constructed Password with the same plaintext verifies too.
	ok, err = mustPassword(t, raw).VerifyPassword(encoded)
	if err != nil || !ok {
		t.Errorf("same plaintext, new instance: ok=%v err=%v", ok, err)
	}
}

func TestVerifyPassword_WrongPassword(t *testing.T) {
	encoded, err := func() (string, error) {
		p := mustPassword(t, "the-right-password")
		return p.HashPassword()
	}()
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	wrongs := []string{
		"the-wrong-password",
		"The-right-password",  // case difference
		"the-right-password ", // trailing space
		"the-right-passwor",
		"the-right-password!",
	}

	for _, w := range wrongs {
		t.Run(w, func(t *testing.T) {
			ok, err := mustPassword(t, w).VerifyPassword(encoded)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok {
				t.Errorf("password %q must not verify", w)
			}
		})
	}
}

func TestVerifyPassword_TamperedHash(t *testing.T) {
	salt := []byte("0123456789abcdef")
	encoded := buildHash(t, "tamper-target", 8, 1, 1, salt, 32)

	parts := strings.Split(encoded, "$")
	raw, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		t.Fatalf("setup: decode hash: %v", err)
	}
	raw[0] ^= 0xFF // flip bits in first byte
	parts[5] = base64.RawStdEncoding.EncodeToString(raw)
	tampered := strings.Join(parts, "$")

	ok, err := mustPassword(t, "tamper-target").VerifyPassword(tampered)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("tampered hash must not verify")
	}
}

func TestVerifyPassword_TamperedSalt(t *testing.T) {
	salt := []byte("0123456789abcdef")
	encoded := buildHash(t, "salt-tamper", 8, 1, 1, salt, 32)

	parts := strings.Split(encoded, "$")
	parts[4] = base64.RawStdEncoding.EncodeToString([]byte("fedcba9876543210"))

	ok, err := mustPassword(t, "salt-tamper").VerifyPassword(strings.Join(parts, "$"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("hash with a different salt must not verify")
	}
}

func TestVerifyPassword_ChangedParamsDoNotVerify(t *testing.T) {
	salt := []byte("0123456789abcdef")
	encoded := buildHash(t, "param-change", 8, 1, 1, salt, 32)

	// Same salt/hash but claim a different iteration count.
	modified := strings.Replace(encoded, "m=8,t=1,p=1", "m=8,t=2,p=1", 1)

	ok, err := mustPassword(t, "param-change").VerifyPassword(modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("altered cost parameters must not verify")
	}
}

// ---------------------------------------------------------------------------
// VerifyPassword – hashes created with non-default parameters
// ---------------------------------------------------------------------------

func TestVerifyPassword_HonoursParamsEncodedInHash(t *testing.T) {
	tests := []struct {
		name        string
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLen     int
		keyLen      uint32
	}{
		{"minimal cost", 8, 1, 1, 16, 32},
		{"short key", 16, 1, 1, 16, 16},
		{"long key", 16, 2, 1, 16, 64},
		{"longer salt", 16, 1, 2, 32, 32},
		{"multiple threads", 32, 2, 4, 16, 32},
	}

	const raw = "portable-hash-test"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt := make([]byte, tt.saltLen)
			for i := range salt {
				salt[i] = byte(i + 1)
			}

			encoded := buildHash(t, raw, tt.memory, tt.iterations, tt.parallelism, salt, tt.keyLen)

			ok, err := mustPassword(t, raw).VerifyPassword(encoded)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Error("expected password to verify against hash with custom params")
			}

			ok, err = mustPassword(t, raw+"x").VerifyPassword(encoded)
			if err != nil {
				t.Fatalf("unexpected error on wrong password: %v", err)
			}
			if ok {
				t.Error("wrong password verified against custom-param hash")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// VerifyPassword – malformed input
// ---------------------------------------------------------------------------

func TestVerifyPassword_InvalidFormat(t *testing.T) {
	validSalt := base64.RawStdEncoding.EncodeToString([]byte("0123456789abcdef"))
	validHash := base64.RawStdEncoding.EncodeToString(make([]byte, 32))

	tests := []struct {
		name    string
		encoded string
		wantErr error
	}{
		{"empty string", "", ErrInvalidHashFormat},
		{"garbage", "not-a-hash", ErrInvalidHashFormat},
		{"too few segments", "$argon2id$v=19$m=8,t=1,p=1$" + validSalt, ErrInvalidHashFormat},
		{
			"too many segments",
			"$argon2id$v=19$m=8,t=1,p=1$" + validSalt + "$" + validHash + "$extra",
			ErrInvalidHashFormat,
		},
		{
			"wrong algorithm argon2i",
			"$argon2i$v=19$m=8,t=1,p=1$" + validSalt + "$" + validHash,
			ErrInvalidHashFormat,
		},
		{
			"wrong algorithm argon2d",
			"$argon2d$v=19$m=8,t=1,p=1$" + validSalt + "$" + validHash,
			ErrInvalidHashFormat,
		},
		{
			"bcrypt style hash",
			"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			ErrInvalidHashFormat,
		},
		{
			"non-numeric version",
			"$argon2id$v=abc$m=8,t=1,p=1$" + validSalt + "$" + validHash,
			ErrIncompatibleVersion,
		},
		{
			"missing version prefix",
			"$argon2id$19$m=8,t=1,p=1$" + validSalt + "$" + validHash,
			ErrIncompatibleVersion,
		},
		{
			"empty version segment",
			"$argon2id$$m=8,t=1,p=1$" + validSalt + "$" + validHash,
			ErrIncompatibleVersion,
		},
		{
			"malformed params",
			"$argon2id$v=19$memory=8$" + validSalt + "$" + validHash,
			ErrInvalidHashFormat,
		},
		{
			"params missing parallelism",
			"$argon2id$v=19$m=8,t=1$" + validSalt + "$" + validHash,
			ErrInvalidHashFormat,
		},
		{
			"params non-numeric",
			"$argon2id$v=19$m=x,t=y,p=z$" + validSalt + "$" + validHash,
			ErrInvalidHashFormat,
		},
		{"empty params", "$argon2id$v=19$$" + validSalt + "$" + validHash, ErrInvalidHashFormat},
		{
			"salt not base64",
			"$argon2id$v=19$m=8,t=1,p=1$!!!notbase64!!!$" + validHash,
			ErrInvalidHashFormat,
		},
		{
			"hash not base64",
			"$argon2id$v=19$m=8,t=1,p=1$" + validSalt + "$!!!notbase64!!!",
			ErrInvalidHashFormat,
		},
		{
			"padded base64 salt rejected",
			"$argon2id$v=19$m=8,t=1,p=1$" + validSalt + "==$" + validHash,
			ErrInvalidHashFormat,
		},
	}

	p := mustPassword(t, "any-valid-password")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := p.VerifyPassword(tt.encoded)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if ok {
				t.Error("ok must be false when an error is returned")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Concurrency / safety
// ---------------------------------------------------------------------------

func TestPassword_ConcurrentHashAndVerify(t *testing.T) {
	p := mustPassword(t, "concurrent-password")

	const workers = 4
	errs := make(chan error, workers)

	for range workers {
		go func() {
			encoded, err := p.HashPassword()
			if err != nil {
				errs <- err
				return
			}
			ok, err := p.VerifyPassword(encoded)
			if err != nil {
				errs <- err
				return
			}
			if !ok {
				errs <- errors.New("verification failed in goroutine")
				return
			}
			errs <- nil
		}()
	}

	for range workers {
		if err := <-errs; err != nil {
			t.Error(err)
		}
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkHashPassword(b *testing.B) {
	p, _ := NewPassword("benchmark-password")

	for b.Loop() {
		if _, err := p.HashPassword(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	p, _ := NewPassword("benchmark-password")
	encoded, err := p.HashPassword()
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		if _, err := p.VerifyPassword(encoded); err != nil {
			b.Fatal(err)
		}
	}
}
