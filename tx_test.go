package microstellar

import "testing"

func TestTx(t *testing.T) {
	ms := New("fake")
	tx := NewTx("fake")

	keyPair, _ := ms.CreateKeyPair()

	err := tx.Sign(keyPair.Seed)

	if err == nil {
		t.Errorf("signing should not succeed: want %v, got nil", err)
	}

	err = tx.Build(sourceAccount(keyPair.Seed))
	if err == nil {
		t.Errorf("build failed: want nil, got %v", err)
	}

	err = tx.Build(sourceAccount(keyPair.Seed))
	if err == nil {
		t.Errorf("duplicate build should fail: want %v, got nil", err)
	}

	tx.Reset()
	err = tx.Build(sourceAccount(keyPair.Seed))
	if err != nil {
		t.Errorf("build failed: want nil, got %v", err)
	}

	err = tx.Sign(keyPair.Seed)
	if err != nil {
		t.Errorf("sign failed: want nil, got %v", err)
	}

	err = tx.Submit()
	if err != nil {
		t.Errorf("submit failed: want nil, got %v", err)
	}

	if tx.Err() != nil {
		t.Errorf("tx.Err() should be nil: got %v", tx.Err())

	}
}
