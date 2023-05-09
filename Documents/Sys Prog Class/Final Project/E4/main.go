//Storing custom data types
package main

import (
    "bytes"
    "encoding/gob"
    "encoding/hex"
    "errors"
    "fmt"
    "log"
    "net/http"
    //"strings"
)

var secret []byte

type User struct {
    Name string
    Age  int
}

func main() {
    gob.Register(&User{})

    var err error
    secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/set", setCookieHandler)
    http.HandleFunc("/get", getCookieHandler)

    log.Print("Listening on :3000")
    err = http.ListenAndServe(":3000", nil)
    if err != nil {
        log.Fatal(err)
    }
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
    user := User{Name: "Alice", Age: 21}

    var buf bytes.Buffer
    err := gob.NewEncoder(&buf).Encode(&user)
    if err != nil {
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }

    cookieValue := buf.Bytes()
    encryptedCookieValue, err := encrypt(cookieValue, secret)
    if err != nil {
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }

    cookie := &http.Cookie{
        Name:     "exampleCookie",
        Value:    hex.EncodeToString(encryptedCookieValue),
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }

    http.SetCookie(w, cookie)
    w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("exampleCookie")
    if err != nil {
        if errors.Is(err, http.ErrNoCookie) {
            http.Error(w, "cookie not found", http.StatusBadRequest)
        } else {
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        return
    }

    encryptedCookieValue, err := hex.DecodeString(cookie.Value)
    if err != nil {
        log.Println(err)
        http.Error(w, "invalid cookie", http.StatusBadRequest)
        return
    }

    cookieValue, err := decrypt(encryptedCookieValue, secret)
    if err != nil {
        log.Println(err)
        http.Error(w, "invalid cookie", http.StatusBadRequest)
        return
    }

    var user User
    err = gob.NewDecoder(bytes.NewReader(cookieValue)).Decode(&user)
    if err != nil {
        log.Println(err)
        http.Error(w, "invalid cookie", http.StatusBadRequest)
        return
    }

    fmt.Fprintf(w, "Name: %q\n", user.Name)
    fmt.Fprintf(w, "Age: %d\n", user.Age)
}

func encrypt(data, key []byte) ([]byte, error) {
    // Implement encryption here
    return data, nil
}

func decrypt(data, key []byte) ([]byte, error) {
    // Implement decryption here
    return data, nil
}
