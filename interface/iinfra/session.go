package iinfra

// ContextKeyTx ...
const ContextKeyTx string = "ContextKeyTx"

type (
	// Tx ...
	Tx interface{}

	// Session ...
	Session interface {
		BeginTx() (Tx, error)
		CommitTx(tx Tx) error
		RollbackTx(tx Tx) error
	}
)
