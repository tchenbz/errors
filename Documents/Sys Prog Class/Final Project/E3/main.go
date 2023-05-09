//Confidential (encrypted) and tamper-proof cookies
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

// customCodec implements the securecookie.Codec interface using a secret key for encryption and decryption
type customCodec struct {
	hashKey  []byte
	blockKey []byte
	secretKey []byte
}

func (c *customCodec) Encode(name string, value interface{}) (string, error) {
	return securecookie.EncodeMulti(name, value, securecookie.New(c.hashKey, c.blockKey))
}

func (c *customCodec) Decode(name string, cookieValue string, value interface{}) error {
	return securecookie.DecodeMulti(name, cookieValue, value, securecookie.New(c.hashKey, c.blockKey))
}

func main() {
	secretKey := []byte("my-secret-key")
	hashKey := securecookie.GenerateRandomKey(32)
	blockKey := securecookie.GenerateRandomKey(16)

	// create a new customCodec instance with the hash key, block key, and secret key
	codec := &customCodec{hashKey, blockKey, secretKey}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		value := map[string]string{
			"username": "sarah",
			"email":    "sarah@example.com",
		}
		// encode the cookie value using the custom codec
		encodedValue, err := codec.Encode("my-cookie-name", value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{
			Name:  "my-cookie-name",
			Value: encodedValue,
			Path:  "/",
			// set HttpOnly and Secure flags to prevent client-side JavaScript from accessing the cookie
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
		fmt.Fprintln(w, "Cookie created successfully.")
	})

	http.HandleFunc("/read-cookie", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("my-cookie-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var value map[string]string
		// decode the cookie value using the custom codec
		err = codec.Decode("my-cookie-name", cookie.Value, &value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Username: %s\nEmail: %s", value["username"], value["email"])
	})

    log.Print("Listening on :4000")
	http.ListenAndServe(":4000", nil)
}

