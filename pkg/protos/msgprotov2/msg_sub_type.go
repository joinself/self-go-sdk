// Copyright 2020 Self Group Ltd. All Rights Reserved.
// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package msgprotov2

import "strconv"

type MsgSubType uint16

const (
	MsgSubTypeUnknown                    MsgSubType = 0
	MsgSubTypeAuthenticationReq          MsgSubType = 1
	MsgSubTypeAuthenticationResp         MsgSubType = 2
	MsgSubTypeAuthenticationQRResp       MsgSubType = 3
	MsgSubTypeAuthenticationDeepLinkResp MsgSubType = 4
	MsgSubTypeFactReq                    MsgSubType = 5
	MsgSubTypeFactResp                   MsgSubType = 6
	MsgSubTypeFactQRResp                 MsgSubType = 7
	MsgSubTypeFactDeepLinkResp           MsgSubType = 8
	MsgSubTypeEmailSecurityCodeReq       MsgSubType = 9
	MsgSubTypeEmailSecurityCodeResp      MsgSubType = 10
	MsgSubTypePhoneSecurityCodeReq       MsgSubType = 11
	MsgSubTypePhoneSecurityCodeResp      MsgSubType = 12
	MsgSubTypePhoneVerificationReq       MsgSubType = 13
	MsgSubTypePhoneVerificationResp      MsgSubType = 14
	MsgSubTypeEmailVerificationReq       MsgSubType = 15
	MsgSubTypeEmailVerificationResp      MsgSubType = 16
	MsgSubTypeDocumentVerificationReq    MsgSubType = 17
	MsgSubTypeDocumentVerificationResp   MsgSubType = 18
)

var EnumNamesMsgSubType = map[MsgSubType]string{
	MsgSubTypeUnknown:                    "Unknown",
	MsgSubTypeAuthenticationReq:          "AuthenticationReq",
	MsgSubTypeAuthenticationResp:         "AuthenticationResp",
	MsgSubTypeAuthenticationQRResp:       "AuthenticationQRResp",
	MsgSubTypeAuthenticationDeepLinkResp: "AuthenticationDeepLinkResp",
	MsgSubTypeFactReq:                    "FactReq",
	MsgSubTypeFactResp:                   "FactResp",
	MsgSubTypeFactQRResp:                 "FactQRResp",
	MsgSubTypeFactDeepLinkResp:           "FactDeepLinkResp",
	MsgSubTypeEmailSecurityCodeReq:       "EmailSecurityCodeReq",
	MsgSubTypeEmailSecurityCodeResp:      "EmailSecurityCodeResp",
	MsgSubTypePhoneSecurityCodeReq:       "PhoneSecurityCodeReq",
	MsgSubTypePhoneSecurityCodeResp:      "PhoneSecurityCodeResp",
	MsgSubTypePhoneVerificationReq:       "PhoneVerificationReq",
	MsgSubTypePhoneVerificationResp:      "PhoneVerificationResp",
	MsgSubTypeEmailVerificationReq:       "EmailVerificationReq",
	MsgSubTypeEmailVerificationResp:      "EmailVerificationResp",
	MsgSubTypeDocumentVerificationReq:    "DocumentVerificationReq",
	MsgSubTypeDocumentVerificationResp:   "DocumentVerificationResp",
}

var EnumValuesMsgSubType = map[string]MsgSubType{
	"Unknown":                    MsgSubTypeUnknown,
	"AuthenticationReq":          MsgSubTypeAuthenticationReq,
	"AuthenticationResp":         MsgSubTypeAuthenticationResp,
	"AuthenticationQRResp":       MsgSubTypeAuthenticationQRResp,
	"AuthenticationDeepLinkResp": MsgSubTypeAuthenticationDeepLinkResp,
	"FactReq":                    MsgSubTypeFactReq,
	"FactResp":                   MsgSubTypeFactResp,
	"FactQRResp":                 MsgSubTypeFactQRResp,
	"FactDeepLinkResp":           MsgSubTypeFactDeepLinkResp,
	"EmailSecurityCodeReq":       MsgSubTypeEmailSecurityCodeReq,
	"EmailSecurityCodeResp":      MsgSubTypeEmailSecurityCodeResp,
	"PhoneSecurityCodeReq":       MsgSubTypePhoneSecurityCodeReq,
	"PhoneSecurityCodeResp":      MsgSubTypePhoneSecurityCodeResp,
	"PhoneVerificationReq":       MsgSubTypePhoneVerificationReq,
	"PhoneVerificationResp":      MsgSubTypePhoneVerificationResp,
	"EmailVerificationReq":       MsgSubTypeEmailVerificationReq,
	"EmailVerificationResp":      MsgSubTypeEmailVerificationResp,
	"DocumentVerificationReq":    MsgSubTypeDocumentVerificationReq,
	"DocumentVerificationResp":   MsgSubTypeDocumentVerificationResp,
}

func (v MsgSubType) String() string {
	if s, ok := EnumNamesMsgSubType[v]; ok {
		return s
	}
	return "MsgSubType(" + strconv.FormatInt(int64(v), 10) + ")"
}
