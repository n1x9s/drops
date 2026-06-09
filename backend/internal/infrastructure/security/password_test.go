package security

import "testing"

func TestPasswordHasher(t *testing.T) {
	hasher := NewPasswordHasher()
	hash, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hash == "correct horse battery staple" {
		t.Fatal("hash should not equal password")
	}
	if !hasher.Compare(hash, "correct horse battery staple") {
		t.Fatal("expected password to match hash")
	}
	if hasher.Compare(hash, "wrong") {
		t.Fatal("wrong password matched hash")
	}
}
