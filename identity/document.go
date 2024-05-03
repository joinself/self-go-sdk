package identity

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl -Wl,--allow-multiple-definition
#cgo darwin LDFLAGS: -lself_sdk
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"

type Document C.self_identity_document

func NewDocument() *Document {
	return nil
}

func (d *Document) NewOperation() *OperationBuilder {
	return &OperationBuilder{}
}
