package security

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	RT = tokenType("rt")
	AT = tokenType("at")
	PT = tokenType("pt")
	UT = tokenType("Unknown")
)

var timeNow func() time.Time = time.Now

type security struct{}

func New() *security {
	return &security{}
}

type tokenType string

func (t tokenType) String() string {
	return strings.ToLower(string(t))
}

func (t tokenType) MarshalJSON() ([]byte, error) {

	switch t.String() {
	case "at", "rt", "pt":
		return json.Marshal(t.String())
	default:
		return nil, fmt.Errorf("error invalid claim value for token type, got: %s", t.String())
	}
}
func newTokenType(tokenType string) tokenType {
	switch strings.ToLower(tokenType) {
	case "rt":
		return RT
	case "at":
		return AT
	case "pt":
		return PT
	default:
		return UT
	}

}

func (t *tokenType) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	vt := newTokenType(v)
	if vt == UT {
		return fmt.Errorf("error invalid claim value for token type, got: %s", v)
	}
	*t = vt
	return nil
}
