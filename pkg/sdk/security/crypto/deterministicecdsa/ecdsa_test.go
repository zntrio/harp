package deterministicecdsa

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func testAllCurves(t *testing.T, f func(*testing.T, elliptic.Curve)) {
	t.Helper()

	tests := []struct {
		name  string
		curve elliptic.Curve
	}{
		{"P256", elliptic.P256()},
		{"P224", elliptic.P224()},
		{"P384", elliptic.P384()},
		{"P521", elliptic.P521()},
	}
	if testing.Short() {
		tests = tests[:1]
	}
	for _, test := range tests {
		curve := test.curve
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			f(t, curve)
		})
	}
}

func TestKeyGeneration(t *testing.T) {
	t.Parallel()

	testAllCurves(t, testKeyGeneration)
}

func testKeyGeneration(t *testing.T, c elliptic.Curve) {
	t.Helper()

	priv, err := GenerateKey(c, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if !c.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		t.Errorf("public key invalid: %s", err)
	}
}

func benchmarkAllCurves(b *testing.B, f func(*testing.B, elliptic.Curve)) {
	tests := []struct {
		name  string
		curve elliptic.Curve
	}{
		{"P256", elliptic.P256()},
		{"P224", elliptic.P224()},
		{"P384", elliptic.P384()},
		{"P521", elliptic.P521()},
	}
	for _, test := range tests {
		curve := test.curve
		b.Run(test.name, func(b *testing.B) {
			f(b, curve)
		})
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	benchmarkAllCurves(b, func(b *testing.B, curve elliptic.Curve) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := GenerateKey(curve, rand.Reader); err != nil {
				b.Fatal(err)
			}
		}
	})
}
