package iinfra

// ContextKeyTx ...
const ContextKeyTx string = "ContextKeyTx"

// Tx ...
type Tx interface{}

// Session ...
type Session interface {
	BeginTx() (Tx, error)
	CommitTx(tx Tx) error
	RollbackTx(tx Tx) error
}
