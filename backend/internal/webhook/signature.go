package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSignature(payload []byte, secret string) string {
	mac := hmac.New(
		sha256.New,
		[]byte(secret),
	)

	mac.Write(payload)

	return hex.EncodeToString(mac.Sum(nil))
}

func VerifySignature(
	payload []byte,
	signature string,
	secret string,
) bool {
	expectedSignature := GenerateSignature(
		payload,
		secret,
	)

	return hmac.Equal(
		[]byte(expectedSignature),
		[]byte(signature),
	)
}

func GenerateTimestampedSignature(
	payload []byte,
	timestamp string,
	secret string,
) string {
	message := timestamp + "." + string(payload)

	mac := hmac.New(
		sha256.New,
		[]byte(secret),
	)

	mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyTimestampedSignature(
	payload []byte,
	timestamp string,
	signature string,
	secret string,
) bool {
	expectedSignature := GenerateTimestampedSignature(
		payload,
		timestamp,
		secret,
	)

	return hmac.Equal(
		[]byte(expectedSignature),
		[]byte(signature),
	)
}
