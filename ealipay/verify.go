package ealipay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

func (c *AlipayClient) VerifySign(params map[string]string, sign string) error {
	if c == nil || c.AlipayPublicKey == nil {
		return fmt.Errorf("alipay public key not configured")
	}
	if sign == "" {
		return fmt.Errorf("missing sign")
	}

	content := buildSignContent(params)
	if content == "" {
		return fmt.Errorf("empty sign content")
	}

	sig, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("decode sign: %w", err)
	}

	sum := sha256.Sum256([]byte(content))
	if err := rsa.VerifyPKCS1v15(c.AlipayPublicKey, crypto.SHA256, sum[:], sig); err != nil {
		return fmt.Errorf("verify sign failed: %w", err)
	}
	return nil
}

func buildSignContent(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "" || v == "" {
			continue
		}
		if k == "sign" || k == "sign_type" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteString("&")
		}
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(params[k])
	}
	return b.String()
}
