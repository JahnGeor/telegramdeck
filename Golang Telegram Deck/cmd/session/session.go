package session

import (
	"GolangTD/utils"
	"errors"
	"net/http"
)

const COOKIE_NAME = "SessionID"
const DONTHAVETHISCOOKIE = "we don't have this cookie in session memory" //код ошибки для отсутствия куки в памяти сессий

var SessionInMemory *Session

type Session struct {
	data map[string]*sessionData //хранит userID в случайно сгенерированных куки
}

type sessionData struct {
	UserID string
}

func NewSession() *Session {
	s := new(Session)
	s.data = make(map[string]*sessionData)
	return s
}

func (s *Session) Init(userID string) string {
	sessionID := utils.GeneratedID()
	s.data[sessionID] = &sessionData{userID}
	return sessionID
}

func (s *Session) Get(sessionID string) *sessionData {
	data := s.data[sessionID]
	if data == nil {
		return nil
	} else {
		return data
	}
}

func CheckCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		return nil, err
	}
	if SessionInMemory.Get(cookie.Value) == nil {
		return nil, errors.New(DONTHAVETHISCOOKIE)
	}
	return cookie, nil
}
