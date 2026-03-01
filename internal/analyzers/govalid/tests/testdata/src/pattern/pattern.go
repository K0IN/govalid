//go:generate govalid ./pattern.go

package pattern

type Pattern struct {
	//govalid:pattern=^[a-z]+$
	Username string `json:"username"`
}
