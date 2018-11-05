package gaekit

import "time"

type TransactionKey_ struct {
	Key       string `datastore:"-" goon:"id"`
	CreatedAt time.Time
}

func NewTransactionKey(key string) *TransactionKey_ {
	return &TransactionKey_{
		Key: key,
	}
}
