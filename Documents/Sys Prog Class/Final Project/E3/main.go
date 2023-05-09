package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/FinalProject/internal/cookies"
)

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
    cookie := http.Cookie{
        Name:     "exampleCookie",
        Value:    "Hello Zoë!",
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }

    err := cookies.WriteEncrypted(w, cookie, secretKey)
    if err != nil {
        log.Println(err)
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }

    w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
    value, err := cookies.ReadEncrypted(r, "exampleCookie", secretKey)
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

    w.Write([]byte(value))
}