package config

type Operation int
type KeyRangeOp int

const (
	READ Operation = iota
	WRITE
)
const (
	TRANSFER KeyRangeOp = iota
	DELETE
)

func (op Operation) String() string {
	return [...]string{"READ", "WRITE"}[op]
}
