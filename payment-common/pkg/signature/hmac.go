package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

func SignHMACSHA256(params map[string]string, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if key == "sign" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys)+1)
	for _, key := range keys {
		parts = append(parts, key+"="+params[key])
	}
	parts = append(parts, "secret="+secret)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strings.Join(parts, "&")))
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyHMACSHA256(params map[string]string, secret string, sign string) bool {
	expected := SignHMACSHA256(params, secret)
	return hmac.Equal([]byte(expected), []byte(sign))
}

