package signature

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestVerification(t *testing.T) {
	type TestCase struct {
		req          *http.Request
		publicKeyPem string
		expected     bool
	}

	testCases := []TestCase{
		{
			req: &http.Request{
				URL:    &url.URL{Path: "/activitypub/users/pichu"},
				Method: "GET",
				Header: map[string][]string{
					"Host":      {"pichuchen.tw"},
					"Accept":    {"application/activity+json, application/ld+json"},
					"Date":      {"Sat, 06 Apr 2024 12:16:05 GMT"},
					"Signature": {"keyId=\"https://g0v.social/actor#main-key\",algorithm=\"rsa-sha256\",headers=\"(request-target) host date accept\",signature=\"F/JFw2aONrK5/UGO5q2M8wIFn/ARGJH90d62Wbpfy1ytGi+qOPM7GGUE8mvWHV7Xo5DHh0vFGUdRlAmqcjjGpZeDSnJtCqdfaOnAd7U0tAq1K2eAR6YTw9E33m6vwGs0+TvsosB3OVcclHnq9Cj+HcDlkDDiuQETDiPLL66Fq3Lfo23yludx6JdmkaP1vc2NkSSxdjUU+uXhjocNQ5lrzAr56PSsoaMnPG/nckyTn1DLJloyBxhlF346M8s+/n1HqrgPGopeEe3eCtbvNz865bCQimodps94QtYQPVppa0IegTzlsuz2bxCW9inHp1cbYE2hXJnNtQsxwucHl95hzA==\""},
				},
			},
			publicKeyPem: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4TyWSRPP1rcsmobFnj1c\n8Qfa1xr9bLH8VMhNF1ClR6aSu0C3f78NWAyCj1w37+2DrS2hPjkjEcyIxl0p3sSW\n7cZ5g16AbVBGUaRLUrUm9lWEK/SXVkXNfVPA67EOWxquP9vxwqOxCRzhE8b8O8qo\nmKtbrGmceokMYfpyromb7knSFsKogcivEZ2FENaxUexYayqxqscwF1JgGGnMx6rh\nWxBwwX40hymgs1vvxOcL1CZKXc1vflzBEZXsw0c67sH2Uk1qxmpNRE8IgYSQ9i5O\nZvF3F6KpCCM6STrCeihJnuZkhedrMfzBxOICfckAXeAvquq6qbqQFwq6v0xp6/11\nuwIDAQAB\n-----END PUBLIC KEY-----\n",
			expected:     true,
		},
		{
			req: &http.Request{
				URL:    &url.URL{Path: "/activitypub/users/pichu"},
				Method: "GET",
				Header: map[string][]string{
					"Host":      {"pichuchen.tw"},
					"Accept":    {`application/activity+json, application/ld+json; profile="https://www.w3.org/ns/activitystreams"`},
					"Date":      {"Sat, 06 Apr 2024 17:58:55 GMT"},
					"Signature": {`keyId="https://misskey.io/users/8xg6igguzx#main-key",algorithm="rsa-sha256",headers="(request-target) date host accept",signature="fkm8WbBeav2qaVMwBtxFkrNnn1eiNGq2jiBfGQESbFBazqVaP22OOT09Dk5XbZqbGE6lTmnbSsjSidKUx45m2u0zQowz97fqROzojGcBT9Ih/SdFFEw7PNqKWo3sPmDL3V/UQQM/mYaZ7x0QsdtOB7iSjOoF83v/E2VCU/Ts5g6hG7XoqoGYDwkpWZXIqzgWloWZgwSFKpZ/FuCM/oshrhFqNXn6G+d5bJeuOtGFnpMf3BqEP+qdFRJ3YxtphV1lcMl3U+iRnlqddzYwZ3G0iuXdFpZlLDYWO9CtZZfPjoUydYXi6Dw2ZuhmPx3YAi3NYs+IMc+Qy6VhbdCpVJ1iEjkzSnfWTNb+wXLbcjy+onCehsr6HEicY0tpprgSwSflUrYHnGI4C5axYgaf/9DfwZ/GVQRwfkD3hr+Whvjgjvs44k2vYjCH/rgc+1v6GmMaSBt3gFshko0oqoIhPYWLBiSXcc/URIkKDmRYaNMEYEgcNUBJSy+H+eI1+2GtAj5JOUM6LfdMhNZdd0i2XIm8+M97aoCtbOvzZEcjYdqW2/hKlYaDSQKOALNeibKnDfgh8/z11bhYZXboq/kHLTiSaVeApJftJ8x/mDC4ZqWwTp6P46wB87XLRGzXyDfZ3oKY68zkaGxTRrw4OYIubK7syr4LU04t+DpjrnX/aMjy9hs="`},
				},
			},
			publicKeyPem: "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAnzXBSB0hQFL4YcWrV00A\nW+Vg0LErgzpqyQSa2KdXMe9KsRojDfCN7jpfGhEG6dILzSmfA/FvXuGkIbiRdmne\n0K/w2uUNbI0gfzxSvS2XMGBeQ927OPWm8LvkqMc5S5JXflReAKjmwHZydxFKIWUH\ns72eX1tcDupgvOZvLJ+ZmDeylATo+n6ZkTgIEPkneCrC15z2rEoDvb+pSEtEyNKy\noXM6+keogz3F7RQ0+IOc+tmPOKXFVYaPiHBuroCidNjNXwZ6N15wTB5lEOrguM28\nutupI8vHH005sxWd9FVyyyrzITpbCyTfP7Y3IwZopS/VJNuBvKQ+oLTM9QQQwse5\nvYNYuCCjMIWm3pEWKCvi7XPEm36ygo5SQRYNh+YOz6+4Qyy+bvO0ETOOvbzL8Af9\njXekilZADMn68E0bGF+iTnkMZw3xd1YUaWAIn//9vAPLNusC+p77xyp2+RKQsw8Q\n/i68oWpmL/gmCpVUpzFnA5TABAZG8m7+bEMht0Ir8uCxF9yL7ilPgTJoBq+wSgIc\nsEqjm1kUSutdQlml+jdc3SqpkBD/bmct+Eli7en7j49uF1KB0Au70pstm30SocMW\nDSOdspRl0UGUuMj2v7szd/FywpvbG2xuzBmJMr/KuYPj3f1wTszD1L2D8+pgcYlJ\nZ+mdv6xNhDPcrgo3r50wJ8UCAwEAAQ==\n-----END PUBLIC KEY-----\n",
			expected:     true,
		},
		{
			req: &http.Request{
				URL:    &url.URL{Path: "/users/pichuchen"},
				Method: "GET",
				Header: map[string][]string{
					"Host":      {"g0v.social"},
					"Date":      {"Sun, 07 Apr 2024 13:04:39 GMT"},
					"Signature": {"keyId=\"https://pichuchen.tw/activitypub/users/pichu#main-key\",algorithm=\"rsa-sha256\",headers=\"(request-target) date host\",signature=\"l5nJhSK+F8kofH/Iak2YcM2JaWrqMyrX6tg7tzaxcwQ+7b+b6C+Pfx757zDKU9StpM0kwgNZXNTQXwYjvsjpbgRbj+1Q6Lnp845yBoVABv/nDlfNe6uYxLt2W9flgB9HH+6EhEkoAURc2eu9t1aH/RzZzUpkfKjXPXf0xtOvCBq99JG9tktn0eJV87MdIeRdwwy4/awyQVdxCBd4c+pCuq4+e2v7g8MSY+URQNSdkApz59e+0HVpjQYb4maHIibYR6QzgxZAHPIaIotXCwLVOxphabl/XcvTHaHmE88Jk/3V7m+DN1RYd0Je9zpDIMN+i+9jd44YulJ4gAOZIaO+2Q==\""},
				},
			},
			publicKeyPem: "-----BEGIN PUBLIC KEY-----\nMIIBCgKCAQEAwWEFx+rcUjw9basgKP2roBPpdvRHks+Rw7bATPz/singTQAm2sMr\ng/DyNJIj2JT6LtiHL91BnQJB79k0okZyvOgFcrZunqOi4cAbKwSRebzOIpP5UE9p\nrjN+25/nUgIUNA631oD6pVu9O4rYAFy8pYL9PGtiXnBtsGbd1Ro1zUA+qgSZpI8A\n1lDrVujGTXQaCMBfgZEvvIeR0AoT3gF2H1X/LkMqMnWRLSK125yz+A0ICUfkimKN\nHCHbpScZ/P114P0WHik3lr6B/oGzDReyXxMLwbnrbulkr1Ne0NgO670IE5sXVq7T\nQDYpzFz+VzmtAZP8SVq2eUzprVd+MUsoiQIDAQAB\n-----END PUBLIC KEY-----\n",
			expected:     true,
		},
	}

	for ti, tc := range testCases {
		actual := VerifySignature(tc.publicKeyPem, tc.req)
		if actual != tc.expected {
			t.Errorf("Test case %d failed: expected %v, got %v", ti, tc.expected, actual)
		}
	}

}

func TestSignatureHeader(t *testing.T) {
	type TestCase struct {
		req           *http.Request
		privateKeyPem string
	}

	testCases := []TestCase{
		{
			req: &http.Request{
				URL:    &url.URL{Path: "/activitypub/users/pichu"},
				Method: "GET",
				Header: map[string][]string{
					"Host":   {"pichuchen.tw"},
					"Accept": {"application/activity+json, application/ld+json"},
					"Date":   {"Sat, 06 Apr 2024 12:16:05 GMT"},
				},
			},
		},
	}

	// Generate PRIVATE KEY
	for ti, _ := range testCases {
		testCases[ti].privateKeyPem = GeneratePrivateKey()
	}

	for ti, tc := range testCases {
		err := Signature(tc.privateKeyPem, "keyId", tc.req)
		if err != nil {
			t.Error(err)
		}
		d, err := httputil.DumpRequest(tc.req, true)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(d))

		publicKeyPem := Pubout([]byte(tc.privateKeyPem))
		actual := VerifySignature(string(publicKeyPem), tc.req)
		if actual != true {
			t.Errorf("Test case %d failed: expected %v, got %v", ti, true, actual)
		}

	}
}

func TestRSAPubout(t *testing.T) {
	privateKeyPem := `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC0S0NTogJhUZ56
3AmWu2BUhxafIVqdG30twyT8zqHsGGPtzITSo23XHZVD9MhEb94dWHv+Io2gXi04
Dq/w8MSudBiFxP7dkwPvH9BVmkIDAT4CXqqldPeE+gwDX3oGM/k3blbERi26EBJd
UDXFuGXXyJLzIXzFJ9xgZ9GohjUlKFJlLXhaErm4Du7BHFsnT02HSXstmfSnIwuX
/DH699FsWGzVDg2Z8ntAuLeVawzfLO6KmxunQLf3U6tX/NrgKrUbUcNuSvTdphFl
UnJDZiMAHd+xkWHez+s0FH3IJA1AXWbjRzsqh2vpzgtRKEVdCLsNo+dB8fdnOVV9
JmnKuLgfAgMBAAECggEAIsIBvE1W7SEZjvD9rj/4bcNPUqVQ/UnP47Mj3dMON2Bq
X21Wy+7y3Y5X+O5nb34rkXe+C7voltqhGBYIyEf6evFpytw0EE5n60E0XlRrVn32
UOmkN1qp30p/Z2UQNsLtUEjm1Fb9OMohaDju7OvEQonp/pJdpfqtyy4oprcc5skg
3tpioZ1HgDMbNFsOtVybI7vCw7KSqU4DU53XzYXSGBSOPyy7yUuPBnAV62vFrXO0
alX3wtCDKMsrO8x+99pV0XtUNPJ5GiAeQiQBFU572fcuAesaOTOzYUG5wtNkgt01
m//QLa/OmnwWhxsK9qUxHpg6UPKohKVhSXqwDVFIWQKBgQDQ1c6fxhvifZi8mKPt
FySKDtSmSUb6XOmOtkv3RlC1+Mgy4pVWRfePd8mBhtEf2dWG30kGZwbnCHp/ok2p
V+NWLq5atbQgmFRDAcX8wCn3b0di28X2wMra1TkUHkP2W5u979KyC1McJaGE6QnS
vO3lybA6tsO6+1vhF/3bC8OCqQKBgQDdA0tS63o/Wui6jWkKJyh0I1Piu06anYc1
yjDx98c4WNW1TCQvFf3kn1L+gjjuBPBGL3tVnhLktm6sV/KhTOzgmDlS2k/E+JKl
mzvAQ06gKhQdcJLarwa0GSyKlpfMSKTBvL/bH0PyrUYlY2/efxMbno5v0zj8Vi/S
K6Ka4fvphwKBgH53fFG1n11hwgU0StFhmncg1BZ/ecRUT/xk2Uw+ANKyfkqSfO1q
Ial4WUdfZXVsbarnvUlrSvrOX6C7oVxUJRAj5OJ590Nl/GE3ybZuiPuXIBv2C8SK
rE4pog9U5MNkZsR3X5kC0qnRACeFGLXORCNMKLTlyVKKtfFyngDw9s1BAoGAJl7W
CvVa6fjsgsbeP6cAvPkNLUX7pZhHyyzpRnkQG1ZA9BLeDVayF6kZjZoqLBirZmQD
859YBGEv4bf7JcnnLi7/dMT2KGpHe3zx6LVGx0PG7j2HIXRVo5rjQsRWYl8SS5hr
bq2E7HLsxLz4xRYmyRrD++Id+KE1+DUfK+ikBa0CgYB9NhLS4bCKlUsGcMfvKe4Y
IGZM5kyAbOkl0oKPQVT/4lsbwExopj9M4oE4BPJcSRx4JWUghjDyHKBvUnMvnPCe
C61EeSrlwKTNRojY09kAKF3FaTaoyJDjeFcWbRJ4oX8QNLPGpzRZFrCCe1WgHRhQ
bgZ5TKo2v+nTWMJq1djPJQ==
-----END PRIVATE KEY-----
`
	expected := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtEtDU6ICYVGeetwJlrtg
VIcWnyFanRt9LcMk/M6h7Bhj7cyE0qNt1x2VQ/TIRG/eHVh7/iKNoF4tOA6v8PDE
rnQYhcT+3ZMD7x/QVZpCAwE+Al6qpXT3hPoMA196BjP5N25WxEYtuhASXVA1xbhl
18iS8yF8xSfcYGfRqIY1JShSZS14WhK5uA7uwRxbJ09Nh0l7LZn0pyMLl/wx+vfR
bFhs1Q4NmfJ7QLi3lWsM3yzuipsbp0C391OrV/za4Cq1G1HDbkr03aYRZVJyQ2Yj
AB3fsZFh3s/rNBR9yCQNQF1m40c7Kodr6c4LUShFXQi7DaPnQfH3ZzlVfSZpyri4
HwIDAQAB
-----END PUBLIC KEY-----
`

	publicKeyPem := Pubout([]byte(privateKeyPem))

	if string(publicKeyPem) != expected {
		t.Errorf("Expected %s, got %s", expected, publicKeyPem)
	}
}
