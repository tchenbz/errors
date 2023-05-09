//Tamper-proof (signed) cookies
package main

import (
    "log"
    "net/http"

    "github.com/gorilla/securecookie"
)

// Declare a global variable to hold the cookie store.
var cookieStore *securecookie.SecureCookie

func main() {
    // Generate a random 32-byte key.
    hashKey := securecookie.GenerateRandomKey(32)
    blockKey := securecookie.GenerateRandomKey(32)

    // Create a new cookie store with the generated keys.
    cookieStore = securecookie.New(hashKey, blockKey)

    mux := http.NewServeMux()
    mux.HandleFunc("/set", setCookieHandler)
    mux.HandleFunc("/get", getCookieHandler)

    log.Print("Listening...")
    err := http.ListenAndServe(":3000", mux)
    if err != nil {
        log.Fatal(err)
    }
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
    // Set a new cookie with the name "exampleCookie" and the value "Hello Zoë!".
    // The Encode() function will sign the cookie value using the keys specified in
    // the cookie store.
    value := map[string]string{
        "message": "Hello Zoë!",
    }
    encoded, err := cookieStore.Encode("exampleCookie", value)
    if err != nil {
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }
    cookie := http.Cookie{
        Name:  "exampleCookie",
        Value: encoded,
        Path:  "/",
    }
    http.SetCookie(w, &cookie)

    w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
    // Get the cookie with the name "exampleCookie". The Decode() function will
    // verify the cookie signature using the keys specified in the cookie store,
    // and return an error if the signature is invalid or if the cookie is not found.
    cookie, err := r.Cookie("exampleCookie")
    if err != nil {
        if err == http.ErrNoCookie {
            http.Error(w, "cookie not found", http.StatusBadRequest)
            return
        }
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }
    value := make(map[string]string)
    err = cookieStore.Decode("exampleCookie", cookie.Value, &value)
    if err != nil {
        http.Error(w, "invalid cookie", http.StatusBadRequest)
        return
    }

    w.Write([]byte(value["message"]))
}

