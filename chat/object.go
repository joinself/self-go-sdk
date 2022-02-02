// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

type Object struct {
	fi         *RemoteFileInteractor
	link       string
	name       string
	mime       string
	expires    int64
	content    string // TODO: this is maybe better a buffer?
	key        string
	nonce      string
	ciphertext string
	logger     string
}

func NewObject(fi *RemoteFileInteractor) *Object {
	return &Object{
		fi: fi,
	}
}

func (o *Object) BuildFromData(name, data, mime string) *Object {

	return o
}
