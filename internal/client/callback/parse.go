package callback

import (
	"errors"
	"net/url"
)

func Parse(fullURL string) (code, state string, err error) {
	u, e := url.Parse(fullURL)
	if e != nil {
		return "", "", e
	}
	q := u.Query()
	code = q.Get("code")
	state = q.Get("state")
	if code == "" || state == "" {
		return "", "", errors.New("missing code or state in callback url")
	}
	return code, state, nil
}
