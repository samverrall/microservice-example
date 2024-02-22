package user

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

const emailMaxLength = 100

var (
	invalidEmailChars = regexp.MustCompile(`[^a-zA-Z0-9+.@_~\-]`)
	validEmailSeq     = regexp.MustCompile(`^[a-zA-Z0-9+._~\-]+@[a-zA-Z0-9+._~\-]+(\.[a-zA-Z0-9+._~\-]+)+$`)
)

type Email string

func NewEmail(email string) (Email, error) {
	if strings.TrimSpace(email) == "" {
		return "", errors.New("cannot be empty")
	}

	if strings.ContainsAny(email, " \t\n\r") {
		return "", errors.New("cannot contain whitespace")
	}
	if strings.ContainsAny(email, `"'`) {
		return "", errors.New("cannot contain quotes")
	}

	if rc := utf8.RuneCountInString(email); rc > emailMaxLength {
		return "", fmt.Errorf("cannot be a over %v characters in length", emailMaxLength)
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		email = strings.TrimSpace(email)
		msg := strings.TrimPrefix(strings.ToLower(err.Error()), "mail: ")

		switch {
		case strings.Contains(msg, "missing '@'"):
			return "", errors.New("missing the @ sign")

		case strings.HasPrefix(email, "@"):
			return "", errors.New("missing part before the @ sign")

		case strings.HasSuffix(email, "@"):
			return "", errors.New("missing part after the @ sign")
		}

		return "", errors.New(msg)
	}

	if addr.Name != "" {
		return "", errors.New("cannot not include a name")
	}

	if matches := invalidEmailChars.FindAllString(addr.Address, -1); len(matches) != 0 {
		return "", fmt.Errorf("cannot contain: %v", matches)
	}

	if !validEmailSeq.MatchString(addr.Address) {
		_, end, _ := strings.Cut(addr.Address, "@")
		if !strings.Contains(end, ".") {
			return "", errors.New("missing top-level domain, e.g. .com, .co.uk, etc.")
		}

		return "", errors.New("must be an email address, e.g. email@example.com")
	}

	return Email(addr.Address), nil
}

func (e Email) String() string {
	return string(e)
}
