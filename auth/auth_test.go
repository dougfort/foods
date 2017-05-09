package auth

import "testing"

func TestAuth(t *testing.T) {
	token1 := []byte("xxx")
	token2 := []byte("yyy")

	// test the basic case: that we get the same string
	a1 := String(token1, "POST", "joe", "banana")
	a2 := String(token1, "POST", "joe", "banana")

	if a1 != a2 {
		t.Fatalf("mismatch: %s != %s", a1, a2)
	}

	// test the get case: with blank food
	a1 = String(token1, "get", "joe", "")
	a2 = String(token1, "GET", "joe", "")

	if a1 != a2 {
		t.Fatalf("mismatch: %s != %s", a1, a2)
	}

	// test that a different token gives a different auth
	a1 = String(token1, "GET", "joe", "")
	a2 = String(token2, "GET", "joe", "")

	if a1 == a2 {
		t.Fatalf("%s == %s", a1, a2)
	}

	// test that a different food gives a different auth
	a1 = String(token1, "POST", "joe", "apple")
	a2 = String(token1, "POST", "joe", "guava")

	if a1 == a2 {
		t.Fatalf("%s == %s", a1, a2)
	}
}
