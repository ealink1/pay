package ealipay

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type AlipayClient struct {
	AppId           string
	PrivateKey      *rsa.PrivateKey
	AlipayPublicKey *rsa.PublicKey
	GatewayUrl      string
}

type Config struct {
	AppId           string
	PrivateKey      string
	AlipayPublicKey string
	IsSandbox       bool
}

func NewClient(config *Config) (*AlipayClient, error) {
	privateKey, err := parsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %w", err)
	}

	alipayPublicKey, err := parsePublicKey(config.AlipayPublicKey)
	if err != nil {
		return nil, fmt.Errorf("parse alipay public key failed: %w", err)
	}

	gatewayUrl := ProdUrl
	if config.IsSandbox {
		gatewayUrl = SandBoxUrl
	}

	return &AlipayClient{
		AppId:           config.AppId,
		PrivateKey:      privateKey,
		AlipayPublicKey: alipayPublicKey,
		GatewayUrl:      gatewayUrl,
	}, nil
}

func parsePrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	privateKeyStr = strings.TrimSpace(privateKeyStr)
	if !strings.HasPrefix(privateKeyStr, "-----BEGIN") {
		privateKeyStr = "-----BEGIN RSA PRIVATE KEY-----\n" + privateKeyStr + "\n-----END RSA PRIVATE KEY-----"
	}

	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return pkcs8Key.(*rsa.PrivateKey), nil
	}
	return key, nil
}

func parsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	publicKeyStr = strings.TrimSpace(publicKeyStr)
	if !strings.HasPrefix(publicKeyStr, "-----BEGIN") {
		publicKeyStr = "-----BEGIN PUBLIC KEY-----\n" + publicKeyStr + "\n-----END PUBLIC KEY-----"
	}

	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub.(*rsa.PublicKey), nil
}

func (c *AlipayClient) sign(data string) (string, error) {
	hashed := crypto.SHA256.New()
	hashed.Write([]byte(data))
	signature, err := rsa.SignPKCS1v15(nil, c.PrivateKey, crypto.SHA256, hashed.Sum(nil))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (c *AlipayClient) buildBizContent(bizContent interface{}) (string, error) {
	if bizContent == nil {
		return "", nil
	}
	data, err := json.Marshal(bizContent)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *AlipayClient) buildCommonParams(method string, bizContent string) map[string]string {
	params := map[string]string{
		"app_id":    c.AppId,
		"method":    method,
		"format":    Format,
		"charset":   Charset,
		"sign_type": SignType,
		"timestamp": "",
		"version":   Version,
	}
	if bizContent != "" {
		params["biz_content"] = bizContent
	}
	return params
}

func (c *AlipayClient) generateSign(params map[string]string) (string, error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(params[k])
	}

	return c.sign(builder.String())
}

func (c *AlipayClient) buildUrl(method string, bizContent interface{}) (string, error) {
	bizContentStr, err := c.buildBizContent(bizContent)
	if err != nil {
		return "", err
	}

	params := c.buildCommonParams(method, bizContentStr)
	params["timestamp"] = getCurrentTimestamp()

	sign, err := c.generateSign(params)
	if err != nil {
		return "", err
	}
	params["sign"] = sign

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return fmt.Sprintf("%s?%s", c.GatewayUrl, values.Encode()), nil
}

func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
