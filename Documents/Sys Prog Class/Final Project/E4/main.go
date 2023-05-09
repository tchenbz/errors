package main

import (
    "bytes"
    "encoding/gob"
    "encoding/hex"
    "errors"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/FinalProject/internal/cookies"
)

var secret []byte

// Declare the User type.
type User struct {
    Name string
    Age  int
}

func main() {
    // Importantly, we need to tell the encoding/gob package about the Go type
    // that we want to encode. We do this my passing *an instance* of the type
    // to gob.Register(). In this case we pass a pointer to an initialized (but
    // empty) instance of the User struct.
    gob.Register(&User{})

    var err error

    secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
    if err != nil {
        log.Fatal(err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/set", setCookieHandler)
    mux.HandleFunc("/get", getCookieHandler)

    log.Print("Listening...")
    err = http.ListenAndServe(":3000", mux)
    if err != nil {
        log.Fatal(err)
    }
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
    // Initialize a User struct containing the data that we want to store in the
    // cookie.
    user := User{Name: "Alice", Age: 21}

    // Initialize a buffer to hold the gob-encoded data.
    var buf bytes.Buffer

    // Gob-encode the user data, storing the encoded output in the buffer.
    err := gob.NewEncoder(&buf).Encode(&user)
    if err != nil {
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }

    // Call buf.String() to get the gob-encoded value as a string and set it as
    // the cookie value.
    cookie := http.Cookie{
        Name:     "exampleCookie",
        Value:    buf.String(),
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }

    // Write an encrypted cookie containing the gob-encoded data as normal.
    err = cookies.WriteEncrypted(w, cookie, secret)
    if err != nil {
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }

    w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
    // Read the gob-encoded value from the encrypted cookie, handling any errors
    // as necessary.
    gobEncodedValue, err := cookies.ReadEncrypted(r, "exampleCookie", secret)
    if err != nil {
        switch {
        case errors.Is(err, http.ErrNoCookie):
            http.Error(w, "cookie not found", http.StatusBadRequest)
        case errors.Is(err, cookies.ErrInvalidValue):
            http.Error(w, "invalid cookie", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        return
    }

    // Create a new instance of a User type.
    var user User

    // Create an strings.Reader containing the gob-encoded value.
    reader := strings.NewReader(gobEncodedValue)

    // Decode it into the User type. Notice that we need to pass a *pointer* to
    // the Decode() target here?
    if err := gob.NewDecoder(reader).Decode(&user); err != nil {
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }

    // Print the user information in the response.
    fmt.Fprintf(w, "Name: %q\n", user.Name)
    fmt.Fprintf(w, "Age: %d\n", user.Age)
}