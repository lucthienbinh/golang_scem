package middlewares

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adam-hanna/sessions"
	"github.com/adam-hanna/sessions/auth"
	"github.com/adam-hanna/sessions/store"
	"github.com/adam-hanna/sessions/transport"
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/models"
)

// SessionJSON is used for marshalling and unmarshalling custom session json information.
// We're using it as an opportunity to tie csrf strings to sessions to prevent csrf attacks
type SessionJSON struct {
	CSRF string `json:"csrf"`
}

var sesh *sessions.Service

// GetUserAuthIDInSession to validate role
func GetUserAuthIDInSession(c *gin.Context) uint {
	r := c.Request
	userSession, _ := sesh.GetUserSession(r)
	rawUint64, _ := strconv.ParseUint(userSession.UserID, 10, 32)
	return uint(rawUint64)
}

// CreateSession after login successful
func CreateSession(c *gin.Context, userAuth *models.UserAuthenticate) {
	w := c.Writer
	csrf, err := generateKey()
	if err != nil {
		log.Printf("Err generating csrf: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}

	myJSON := SessionJSON{
		CSRF: csrf,
	}
	JSONBytes, err := json.Marshal(myJSON)
	if err != nil {
		log.Printf("Err marhsalling json: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}

	userID := strconv.FormatUint(uint64(userAuth.ID), 10)
	userSession, err := sesh.IssueUserSession(userID, string(JSONBytes[:]), w)
	if err != nil {
		log.Printf("Err issuing user session: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}
	log.Printf("In issue; user's session: %v\n", userSession)

	// note: we set the csrf in a cookie, but look for it in request headers
	csrfCookie := http.Cookie{
		Name:     "csrf",
		Value:    csrf,
		Expires:  userSession.ExpiresAt,
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // note: can't use secure cookies in development
	}
	http.SetCookie(w, &csrfCookie)

	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}

// ValidateSession when user access route with middleware
func ValidateSession(c *gin.Context) {
	w := c.Writer
	r := c.Request
	userSession, err := sesh.GetUserSession(r)
	if err != nil {
		log.Printf("Err fetching user session: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}
	// nil session pointers indicate a 401 unauthorized
	if userSession == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	log.Printf("In validate; user session expiration before extension: %v\n", userSession.ExpiresAt.UTC())

	myJSON := SessionJSON{}
	if err := json.Unmarshal([]byte(userSession.JSON), &myJSON); err != nil {
		log.Printf("Err unmarshalling json: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}
	log.Printf("In validate; user's custom json: %v\n", myJSON)

	// note: we set the csrf in a cookie, but look for it in request headers
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != myJSON.CSRF {
		log.Printf("Unauthorized! CSRF token doesn't match user session")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	// note that session expiry's need to be manually extended
	if err = sesh.ExtendUserSession(userSession, r, w); err != nil {
		log.Printf("Err extending user session: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}
	log.Printf("In validate; users session expiration after extension: %v\n", userSession.ExpiresAt.UTC())

	// need to extend the csrf cookie, too
	csrfCookie := http.Cookie{
		Name:     "csrf",
		Value:    csrf,
		Expires:  userSession.ExpiresAt,
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // note: can't use secure cookies in development
	}
	http.SetCookie(w, &csrfCookie)
	c.Next()
}

// ClearSession when user loged out
func ClearSession(c *gin.Context) {
	w := c.Writer
	r := c.Request
	userSession, err := sesh.GetUserSession(r)
	if err != nil {
		log.Printf("Err fetching user session: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}
	// nil session pointers indicate a 401 unauthorized
	if userSession == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	log.Printf("In clear; session: %v\n", userSession)

	myJSON := SessionJSON{}
	if err := json.Unmarshal([]byte(userSession.JSON), &myJSON); err != nil {
		log.Printf("Err unmarshalling json: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}
	log.Printf("In clear; user's custom json: %v\n", myJSON)

	// note: we set the csrf in a cookie, but look for it in request headers
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != myJSON.CSRF {
		log.Printf("Unauthorized! CSRF token doesn't match user session")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	if err = sesh.ClearUserSession(userSession, w); err != nil {
		log.Printf("Err clearing user session: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error"})
		return
	}

	// need to clear the csrf cookie, too
	aLongTimeAgo := time.Now().Add(-1000 * time.Hour)
	csrfCookie := http.Cookie{
		Name:     "csrf",
		Value:    "",
		Expires:  aLongTimeAgo,
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // note: can't use secure cookies in development
	}
	http.SetCookie(w, &csrfCookie)

	// w.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, gin.H{"status": "You are logged out"})
}

// RunAuth to connect redis server
func RunAuth() {
	seshStore := store.New(store.Options{})

	// e.g. `$ openssl rand -base64 64`
	seshAuth, err := auth.New(auth.Options{
		Key: []byte("qV3Dyaji0KpAgO5yTk3tskcwvLFJZk7LTzcptpbc2ecFHI9kM763UTXyexZp/DcXLBr15hz16PjrY6qjn8Xmvw=="),
	})
	if err != nil {
		log.Fatal(err)
	}

	seshTransport := transport.New(transport.Options{
		HTTPOnly: true,
		Secure:   false, // note: can't use secure cookies in development!
	})

	sesh = sessions.New(seshStore, seshAuth, seshTransport, sessions.Options{})

}

// thanks
// https://astaxie.gitbooks.io/build-web-application-with-golang/en/06.2.html#unique-session-ids
func generateKey() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
