package user

type Password string

func NewPassword(p string) (Password, error) {
	return "", nil
}

func (p Password) String() string {
	return string(p)
}
