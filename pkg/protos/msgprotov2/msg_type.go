// Copyright 2020 Self Group Ltd. All Rights Reserved.
// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package msgprotov2

import "strconv"

type MsgType int8

const (
	MsgTypeMSG  MsgType = 0
	MsgTypeACK  MsgType = 1
	MsgTypeERR  MsgType = 2
	MsgTypeAUTH MsgType = 3
	MsgTypeWTC  MsgType = 5
	MsgTypeSTS  MsgType = 6
)

var EnumNamesMsgType = map[MsgType]string{
	MsgTypeMSG:  "MSG",
	MsgTypeACK:  "ACK",
	MsgTypeERR:  "ERR",
	MsgTypeAUTH: "AUTH",
	MsgTypeWTC:  "WTC",
	MsgTypeSTS:  "STS",
}

var EnumValuesMsgType = map[string]MsgType{
	"MSG":  MsgTypeMSG,
	"ACK":  MsgTypeACK,
	"ERR":  MsgTypeERR,
	"AUTH": MsgTypeAUTH,
	"WTC":  MsgTypeWTC,
	"STS":  MsgTypeSTS,
}

func (v MsgType) String() string {
	if s, ok := EnumNamesMsgType[v]; ok {
		return s
	}
	return "MsgType(" + strconv.FormatInt(int64(v), 10) + ")"
}
