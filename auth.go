
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

var (
	// Options
	flagAlg     = jwtCmd.Flags().String("alg", "RS256", "signing algorithm identifier")
	flagCompact = jwtCmd.Flags().Bool("compact", false, "output compact JSON")
	flagDebug   = jwtCmd.Flags().Bool("debug", false, "print out all kinds of debug data")

	// Modes - exactly one of these is required
	flagSign   = jwtCmd.Flags().String("sign", "", "path to claims object to sign or '-' to read from stdin")
	flagVerify = jwtCmd.Flags().String("verify", "", "path to JWT token to verify or '-' to read from stdin")
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "generated a jwt token",
	Long:  "echo {\"foo\":\"bar\"} | bin/jwt -key test/sample_key -alg RS256 -sign ",
	//RunE:
}

func init() {
	jwtCmd.RunE = runJwt
}

func runJwt(cmd *cobra.Command, args []string) error {
	if *flagSign != "" {
		return signToken()
	} else if *flagVerify != "" {
		return verifyToken()
	} else {
		flag.Usage()
		return fmt.Errorf("None of the required flags are present.  What do you want me to do?")
	}
}

// A useful example app.  You can use this to debug your tokens on the command line.
// This is also a great place to look at how you might use this library.
//
// Example usage:
// The following will create and sign a token, then verify it and output the original claims.
//     echo {\"foo\":\"bar\"} | bin/jwt -key test/sample_key -alg RS256 -sign - | bin/jwt -key test/sample_key.pub -verify -

// Helper func:  Read input from specified file or stdin
func loadData(p string) ([]byte, error) {
	return []byte(p), nil
}

// Print a json object in accordance with the prophecy (or the command line options)
func printJSON(j interface{}) error {
	var out []byte
	var err error

	if *flagCompact == false {
		out, err = json.MarshalIndent(j, "", "    ")
	} else {
		out, err = json.Marshal(j)
	}

	if err == nil {
		fmt.Println(string(out))
	}

	return err
}

// Verify a token and output the claims.  This is a great example
// of how to verify and view a token.
func verifyToken() error {
	// get the token
	tokData, err := loadData(*flagVerify)
	if err != nil {
		return fmt.Errorf("Couldn't read token: %v", err)
	}

	// trim possible whitespace from token
	tokData = regexp.MustCompile(`\s*$`).ReplaceAll(tokData, []byte{})
	if *flagDebug {
		fmt.Fprintf(os.Stderr, "Token len: %v bytes\n", len(tokData))
	}

	// Parse the token.  Load the key from command line option
	token, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		data, err := loadData(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4Q+gobOvKxGZV1dqIOrMenSuBWmx4byLjl2+R15LCsgr12YK
iVqf9PWqMMXXn09CqSSOSKgl26epmIgJsg7ERHpbTEtr9tkusKNin196O4m+I+MA
vISgI28QzfHSjysz4DOA7KI49CAUVW2aY1MHXhUlYb9CwUka0CIwhSzxARsH0YDJ
e5aJnd7AqxML9ewy8cn8S7c+24/j67Hs73vghobJvytW2IMkUDPDdbMwOu+4q6D1
sQSMMaRtdVm1l7+62viDQ0ruDEHarRbtYo24UH/NpR5LX/9WeWRW8b742Cb2kgY+
tRMUug1akN8OXPXGqQMvNRwf2b/GkmZgEqNvCQIDAQABAoIBACPb45IgGQbZtP7v
lJ9OCepw8NE39/mpmurCEPE6GubK4gFd5svfvqC/C7DdtO9TJ6HfizQUZoPLAQCm
nDTcmXT3sdhEJB2emQvX3HzcL5OQ7NS29IUU9JbwuVK29v+MuKU+T4pkhoKPIe5C
mli8/+2DnQMoADDfvv/ukqCFepjcMzNBsdhXYAboHv6U3j5bs27SGuqXU5NEzfuX
clOFjTTsM3s4wqSti8ZpUWZ5rN7+FB4O3jNmMXQr/gF+ceeGOpXIltT3JtVHHRI5
f/v2cGOfQl3k0upeafWyElIfblgKCQfKAcdiyo95NybrNmRQ7Xkx9nBvcNLHO2Bx
daqdTAECgYEA8IcwvMZwHvkYQABL/9rBntnQW+Jgn99a5lvlfjjaW/9w3sA2HXI6
rJhDeG8X4p5gg4BegNOeMST1L3sfiTw5/DTVoymqmA3FPkciW9obyH5ZJWIivq/F
3Xv1wyuIxPvChl/78maL2+dEsSTkJ4RvS80p997DfmU1KUbSRobBL8ECgYEA74m9
N/QgkpkzNeyIqEtd68BdaioH9+p/MgRdCDYR7pLqwlynUIH7YIbtmdOWEpvZg543
nSnIf5HrjEucTXSTLVDKWy7oLIB5UTLAS2FbB6Y57WspCDi/hSW8cNuxtGOEbjB+
qdO8eQy+/K0WLkxO7lPxi51lyiWleweAJAmhEUkCgYBPoWJITSYXiv41SiPfG8xY
S+JIWUUGCMsuUqRCyo24QXRbuqTv0L6OH4bO23C77RUk1B31ZpoLySGHS6rgI1lL
Hy7Pat74oi539NLyN95U3UekMb4xBT5rmjt+Fu6b0IHRPPvLf5mz/vfl8cG7N4Ql
Q1IupshwEw+rj6/T+47/wQKBgQCIrNzsWj1jqEpSEF6BOE+kvqQOeWEGkiR1U4wJ
rWBZ8jZFJDzLcP8Puq1Dwji08XwQ32v4Hukp8QanjFTo1QVNK/XqRT9wdPXD4ONb
n3cjTDNtRmGpMUgGHtwAwToKJWZgwQbku82kfCNVZSVs0VmQHxGJiguUZhqfsk3p
Qh1HEQKBgGvmVXPNlEohvZ3O5Yaqgg1tYq9TaQZjjobOPg9GGm7yG81epfH0jljt
jbsRkgCnYCYyIZ7fo1ro/7GKN3lbmwNxF0w+qLvknW9NdpBprpe413nfmGe4n5oQ
aRFIFgBBn/KRM2/iwl3tg7eORscKj9E8kdu/QQtLVgA9qz2daB8t
-----END RSA PRIVATE KEY----`)
		if err != nil {
			return nil, err
		}
		if isEs() {
			return jwt.ParseECPublicKeyFromPEM(data)
		}
		return data, nil
	})

	// Print some debug data
	if *flagDebug && token != nil {
		fmt.Fprintf(os.Stderr, "Header:\n%v\n", token.Header)
		fmt.Fprintf(os.Stderr, "Claims:\n%v\n", token.Claims)
	}

	// Print an error if we can't parse for some reason
	if err != nil {
		return fmt.Errorf("Couldn't parse token: %v", err)
	}

	// Is token invalid?
	if !token.Valid {
		return fmt.Errorf("Token is invalid")
	}

	// Print the token details
	if err := printJSON(token.Claims); err != nil {
		return fmt.Errorf("Failed to output claims: %v", err)
	}

	return nil
}

// Create, sign, and output a token.  This is a great, simple example of
// how to use this library to create and sign a token.
func signToken() error {
	// get the token data from command line arguments
	tokData, err := loadData(*flagSign)
	if err != nil {
		return fmt.Errorf("Couldn't read token: %v", err)
	} else if *flagDebug {
		fmt.Fprintf(os.Stderr, "Token: %v bytes", len(tokData))
	}

	// parse the JSON of the claims
	var claims map[string]interface{}
	fmt.Println(tokData)
	if err := json.Unmarshal(tokData, &claims); err != nil {
		return fmt.Errorf("Couldn't parse claims JSON: %v", err)
	}

	// get the key
	var key interface{}
	key, err = loadData(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4Q+gobOvKxGZV1dqIOrMenSuBWmx4byLjl2+R15LCsgr12YK
iVqf9PWqMMXXn09CqSSOSKgl26epmIgJsg7ERHpbTEtr9tkusKNin196O4m+I+MA
vISgI28QzfHSjysz4DOA7KI49CAUVW2aY1MHXhUlYb9CwUka0CIwhSzxARsH0YDJ
e5aJnd7AqxML9ewy8cn8S7c+24/j67Hs73vghobJvytW2IMkUDPDdbMwOu+4q6D1
sQSMMaRtdVm1l7+62viDQ0ruDEHarRbtYo24UH/NpR5LX/9WeWRW8b742Cb2kgY+
tRMUug1akN8OXPXGqQMvNRwf2b/GkmZgEqNvCQIDAQABAoIBACPb45IgGQbZtP7v
lJ9OCepw8NE39/mpmurCEPE6GubK4gFd5svfvqC/C7DdtO9TJ6HfizQUZoPLAQCm
nDTcmXT3sdhEJB2emQvX3HzcL5OQ7NS29IUU9JbwuVK29v+MuKU+T4pkhoKPIe5C
mli8/+2DnQMoADDfvv/ukqCFepjcMzNBsdhXYAboHv6U3j5bs27SGuqXU5NEzfuX
clOFjTTsM3s4wqSti8ZpUWZ5rN7+FB4O3jNmMXQr/gF+ceeGOpXIltT3JtVHHRI5
f/v2cGOfQl3k0upeafWyElIfblgKCQfKAcdiyo95NybrNmRQ7Xkx9nBvcNLHO2Bx
daqdTAECgYEA8IcwvMZwHvkYQABL/9rBntnQW+Jgn99a5lvlfjjaW/9w3sA2HXI6
rJhDeG8X4p5gg4BegNOeMST1L3sfiTw5/DTVoymqmA3FPkciW9obyH5ZJWIivq/F
3Xv1wyuIxPvChl/78maL2+dEsSTkJ4RvS80p997DfmU1KUbSRobBL8ECgYEA74m9
N/QgkpkzNeyIqEtd68BdaioH9+p/MgRdCDYR7pLqwlynUIH7YIbtmdOWEpvZg543
nSnIf5HrjEucTXSTLVDKWy7oLIB5UTLAS2FbB6Y57WspCDi/hSW8cNuxtGOEbjB+
qdO8eQy+/K0WLkxO7lPxi51lyiWleweAJAmhEUkCgYBPoWJITSYXiv41SiPfG8xY
S+JIWUUGCMsuUqRCyo24QXRbuqTv0L6OH4bO23C77RUk1B31ZpoLySGHS6rgI1lL
Hy7Pat74oi539NLyN95U3UekMb4xBT5rmjt+Fu6b0IHRPPvLf5mz/vfl8cG7N4Ql
Q1IupshwEw+rj6/T+47/wQKBgQCIrNzsWj1jqEpSEF6BOE+kvqQOeWEGkiR1U4wJ
rWBZ8jZFJDzLcP8Puq1Dwji08XwQ32v4Hukp8QanjFTo1QVNK/XqRT9wdPXD4ONb
n3cjTDNtRmGpMUgGHtwAwToKJWZgwQbku82kfCNVZSVs0VmQHxGJiguUZhqfsk3p
Qh1HEQKBgGvmVXPNlEohvZ3O5Yaqgg1tYq9TaQZjjobOPg9GGm7yG81epfH0jljt
jbsRkgCnYCYyIZ7fo1ro/7GKN3lbmwNxF0w+qLvknW9NdpBprpe413nfmGe4n5oQ
aRFIFgBBn/KRM2/iwl3tg7eORscKj9E8kdu/QQtLVgA9qz2daB8t
-----END RSA PRIVATE KEY----`)
	if err != nil {
		return fmt.Errorf("Couldn't read key: %v", err)
	}

	// get the signing alg
	alg := jwt.GetSigningMethod(*flagAlg)
	if alg == nil {
		return fmt.Errorf("Couldn't find signing method: %v", *flagAlg)
	}

	// create a new token
	token := jwt.New(alg)
	token.Claims = claims

	if isEs() {
		if k, ok := key.([]byte); !ok {
			return fmt.Errorf("Couldn't convert key data to key")
		} else {
			key, err = jwt.ParseECPrivateKeyFromPEM(k)
			if err != nil {
				return err
			}
		}
	}

	if out, err := token.SignedString(key); err == nil {
		fmt.Println(out)
	} else {
		return fmt.Errorf("Error signing token: %v", err)
	}

	return nil
}

func isEs() bool {
	return strings.HasPrefix(*flagAlg, "ES")
}
