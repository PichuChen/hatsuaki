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
