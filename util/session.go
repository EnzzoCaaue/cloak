package util

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	store = sessions.NewCookieStore([]byte("1234567891234567"))
)

// Session is the session struct
type Session struct {
	store *sessions.Session
}

// GetSession gets a request session
func GetSession(req *http.Request, name string) (*Session, error) {
	sess, err := store.Get(req, name)
	if err != nil {
		return nil, err
	}
	s := &Session{
		sess,
	}
	return s, err
}

// GetInt gets a int value from a session
func (s *Session) GetInt(key string) int {
	if i, ok := s.store.Values[key].(int); ok {
		return i
	}
	return 0
}

// GetString gets a string value from a session
func (s *Session) GetString(key string) string {
	if i, ok := s.store.Values[key].(string); ok {
		return i
	}
	return ""
}

// Set sets a session value
func (s *Session) Set(key string, val interface{}) {
	s.store.Values[key] = val
}

// Delete deletes a session value
func (s *Session) Delete(key string) {
	delete(s.store.Values, key)
	s.store.Options.MaxAge = -1
}

// AddFlash adds a flash string to the session
func (s *Session) AddFlash(msg, key string) {
	s.store.AddFlash(msg, key)
}

// GetFlashes returns an array of string flashes
func (s *Session) GetFlashes(key string) []string {
	flashes := s.store.Flashes(key)
	msgs := []string{}
	for i := range flashes {
		if v, ok := flashes[i].(string); ok {
			msgs = append(msgs, v)
		}
	}
	return msgs
}

// Save saves the session to the request
func (s *Session) Save(req *http.Request, w http.ResponseWriter) error {
	err := s.store.Save(req, w)
	return err
}
