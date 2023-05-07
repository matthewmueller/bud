package session

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/matthewmueller/bud/cipher"
	"github.com/matthewmueller/bud/middleware"
)

type Middleware = middleware.Middleware

const sessionID = "sid"

func New[Session any](cipher cipher.Cipher, store Store) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionID)
			if err != nil {
				if !errors.Is(err, http.ErrNoCookie) {
					fmt.Println("cookie error, showing unauthorized", err)
					http.Error(w, "cookie get error", http.StatusUnauthorized)
					return
				}
				cookie = &http.Cookie{
					Name:     sessionID,
					Value:    "some-uuid",
					Path:     "/",
					HttpOnly: true,
					Expires:  time.Now().Add(24 * time.Hour),
				}
				session := new(Session)
				_ = session
				// TODO: short-circuit because we wouldn't have a session if we don't
				// have a cookie
			}
			plainID, err := cipher.Decrypt([]byte(cookie.Value))
			if err != nil {
				fmt.Println("cookie decrypt error, showing unauthorized", err)
				http.Error(w, "cookie decrypt error", http.StatusUnauthorized)
				return
			}
			sessionRaw, err := store.Get(string(plainID))
			if err != nil {
				if !errors.Is(err, ErrNotFound) {
					fmt.Println("session get error, showing unauthorized", err)
					http.Error(w, "session get error", http.StatusUnauthorized)
					return
				}
				// session = NewSession()
				_ = sessionRaw
			}
			// fmt.Println("got c", )
		})
	}
}

// func Middleware(secret string) middleware.Middleware {
// 	// cipher := secretbox.New([32]byte{
// 	// 	0xf5, 0xaf, 0xe2, 0xcb, 0x87, 0xfb, 0x59, 0x65, 0x3d, 0xff,
// 	// 	0x43, 0x56, 0x19, 0x4a, 0x22, 0x64, 0x91, 0x4a, 0x28, 0xa0,
// 	// 	0x4a, 0x06, 0xb8, 0x21, 0x29, 0x42, 0xb4, 0x44, 0x55, 0xd1,
// 	// 	0x13, 0x89,
// 	// })
// 	// cs := cookies.Secure(cipher)
// 	// return &Middleware{
// 	// 	Cookies: cookies.Default(),
// 	// 	Store:   nil,
// 	// }

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// type Session[Data any] struct {
// 	ID   string
// 	Data Data
// }

// type Middleware struct {
// 	Cookies cookies.Store
// 	Store   budsession.Store // If nil, we use the cookie store
// }

// const sessionID = "sid"

// type contextKey string

// var sessionKey = contextKey("session")

// func (m *Middleware) Middleware(next http.Handler) http.Handler {
// 	cipher := secretbox.New([32]byte{
// 		0xf5, 0xaf, 0xe2, 0xcb, 0x87, 0xfb, 0x59, 0x65, 0x3d, 0xff,
// 		0x43, 0x56, 0x19, 0x4a, 0x22, 0x64, 0x91, 0x4a, 0x28, 0xa0,
// 		0x4a, 0x06, 0xb8, 0x21, 0x29, 0x42, 0xb4, 0x44, 0x55, 0xd1,
// 		0x13, 0x89,
// 	})
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		store := m.Store
// 		if store == nil {
// 			store = budsession.New(cookies.Secure(cipher), w, r)
// 		}

// 		cookie, err := m.Cookies.Get(r, sessionID)
// 		if err != nil {
// 			if !errors.Is(err, http.ErrNoCookie) {
// 				fmt.Println("cookie error, showing unauthorized", err)
// 				http.Error(w, "cookie get error", http.StatusUnauthorized)
// 				return
// 			}
// 			cookie = &http.Cookie{
// 				Name:     sessionID,
// 				Value:    "a123",
// 				Path:     "/",
// 				HttpOnly: true,
// 				MaxAge:   86400, // 1 day
// 				Expires:  time.Now().Add(86400 * time.Second),
// 			}
// 		}

// 		var session Session
// 		if sessionData, err := store.Get(cookie.Value); err != nil {
// 			if !errors.Is(err, budsession.ErrNotFound) {
// 				fmt.Println("session get error, showing unauthorized", err)
// 				http.Error(w, "session get error", http.StatusUnauthorized)
// 				return
// 			}
// 		} else if sessionData != nil {
// 			fmt.Println("GOT SESSION DATA", string(sessionData))
// 			if err := json.NewDecoder(bytes.NewReader(sessionData)).Decode(&session); err != nil {
// 				fmt.Println("session decode error, showing unauthorized", err)
// 				http.Error(w, "session decode error", http.StatusUnauthorized)
// 				return
// 			}
// 		}

// 		ctx := context.WithValue(r.Context(), sessionKey, &session)
// 		r = r.WithContext(ctx)

// 		fmt.Println("mw: session before", session.UserID)
// 		next.ServeHTTP(w, r)

// 		payload := new(bytes.Buffer)
// 		if err := json.NewEncoder(payload).Encode(session); err != nil {
// 			fmt.Println("session encode error, showing unauthorized", err)
// 			http.Error(w, "session encode error", http.StatusUnauthorized)
// 			return
// 		}
// 		fmt.Println(payload.String())

// 		if err := store.Set(cookie.Value, payload.Bytes(), time.Now()); err != nil {
// 			fmt.Println("session set error, showing unauthorized", err)
// 			http.Error(w, "session set error", http.StatusUnauthorized)
// 			return
// 		}

// 		if err := m.Cookies.Set(w, cookie); err != nil {
// 			fmt.Println("cookie set error, showing unauthorized", err)
// 			http.Error(w, "cookie set error", http.StatusUnauthorized)
// 			return
// 		}
// 	})
// }
