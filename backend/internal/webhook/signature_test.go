package webhook

import "testing"

func TestGenerateSignature(t *testing.T) {

	payload := []byte(`{"event_id":"evt-001","amount":50000}`)
	secret := "webhook-secret"

	signature := GenerateSignature(
		payload,
		secret,
	)

	if signature == "" {
		t.Fatal("signature should not be empty")
	}
}

func TestVerifySignature(t *testing.T) {

	payload := []byte(`{"event_id":"evt-001","amount":50000}`)
	secret := "webhook-secret"

	signature := GenerateSignature(
		payload,
		secret,
	)

	if !VerifySignature(
		payload,
		signature,
		secret,
	) {
		t.Fatal("signature should be valid")
	}
}

func TestVerifySignatureInvalid(t *testing.T) {

	payload := []byte(`{"event_id":"evt-001","amount":50000}`)
	secret := "webhook-secret"

	signature := GenerateSignature(
		payload,
		secret,
	)

	if VerifySignature(
		payload,
		signature+"invalid",
		secret,
	) {
		t.Fatal("signature should be invalid")
	}
}
