package token

type Token string

func (this_ Token) String() string {
	if this_ == "" {
		return "UNKNOWN"
	}
	return "token(" + string(this_) + ")"
}
