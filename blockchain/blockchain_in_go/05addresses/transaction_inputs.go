package main

import (
	"bytes"
)


// TXInput represents a transaction input.
type TXInput struct {
	TXid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}


// UsesKey checks whether the addresses initiated the transaction.
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
