// Copyright 2020 Self Group Ltd. All Rights Reserved.

package chat

import (
	"log"
	"strconv"
)

type Object struct {
	fi         *RemoteFileInteractor
	Link       string
	Name       string
	Mime       string
	Expires    int64
	Key        string
	Nonce      string
	Ciphertext string
}

func NewObject(fi *RemoteFileInteractor) *Object {
	return &Object{
		fi: fi,
	}
}

func (o *Object) BuildFromData(data []byte, name, mime string) error {
	println("building an object")
	// Encrypt message
	obj, err := o.fi.SetObject(data)
	if err != nil {
		log.Println("error uploading the object")
		return err
	}

	o.Name = name
	o.Link = obj.Link
	o.Key = obj.Key
	o.Mime = mime
	o.Expires = obj.Expires

	return nil
}

func (o *Object) BuildFromObject(obj map[string]interface{}) error {
	o.Name = obj["name"].(string)
	o.Link = obj["link"].(string)
	if _, ok := obj["key"]; ok {
		o.Key = obj["key"].(string)
		o.Mime = obj["mime"].(string)
		o.Expires, _ = strconv.ParseInt(obj["expires"].(string), 10, 64)
	}

	return nil
}

func (o *Object) ToPayload() map[string]interface{} {
	return map[string]interface{}{
		"name":    o.Name,
		"link":    o.Link,
		"key":     o.Key,
		"mime":    o.Mime,
		"expires": strconv.FormatInt(o.Expires, 10),
	}
}

func (o *Object) GetContent() ([]byte, error) {
	content, err := o.fi.GetObject(o.Link, o.Key)
	if err != nil {
		println("error getting the object")
		println(err.Error())
		return []byte(""), err
	}

	return content, nil
}
