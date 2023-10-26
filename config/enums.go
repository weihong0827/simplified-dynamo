package config

type Operation int

const (
	READ Operation = iota
	WRITE
)

func (op Operation) String() string {
	return [...]string{"READ", "WRITE"}[op]
}
