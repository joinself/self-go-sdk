// Copyright 2020 Self Group Ltd. All Rights Reserved.

package fact

import (
	"encoding/json"
	"errors"

	"github.com/square/go-jose"

	"github.com/tidwall/gjson"
)

var (
	OperatorEqual              = "=="
	OperatorDifferent          = "!="
	OperatorGreaterOrEqualThan = ">="
	OperatorLessOrEqualThan    = "<="
	OperatorGreaterThan        = ">"
	OperatorLessThan           = "<"

	RequestInformation  = "identities.facts.query.req"
	ResponseInformation = "identities.facts.query.resp"

	StatusAccepted     = "accepted"
	StatusRejected     = "rejected"
	StatusUnauthorized = "unauthorized"

	ErrFactEmptyName           = errors.New("provided fact does not specify a name")
	ErrFactBadSource           = errors.New("fact must specify one source")
	ErrFactInvalidSource       = errors.New("provided fact does not specify a valid source")
	ErrFactInvalidOperator     = errors.New("provided fact does not specify a valid operator")
	ErrFactInvalidFactToSource = errors.New("provided source does not support given fact")
	ErrInvalidSourceSpec       = errors.New("internal error : invalid source spec")
)

// Fact specific details about the fact
type Fact struct {
	Fact          string            `json:"fact"`
	Sources       []string          `json:"sources,omitempty"`
	Origin        string            `json:"iss,omitempty"`
	Operator      string            `json:"operator,omitempty"`
	Attestations  []json.RawMessage `json:"attestations,omitempty"`
	Issuers       []string          `json:"issuers,omitempty"`
	ExpectedValue string            `json:"expected_value,omitempty"`
	AttestedValue string            `json:"-"`
	payloads      [][]byte
	results       []string
	value         string
}

// AttestedValues returns all attested values for an attestations
func (f *Fact) AttestedValues() []string {
	values := make([]string, len(f.payloads))

	for i, p := range f.payloads {
		v := gjson.GetBytes(p, f.Fact).String()
		if v == "" {
			for _, sf := range gjson.GetBytes(p, "facts").Array() {
				if sf.Map()["key"].String() == f.Fact {
					v = sf.Map()["value"].String()
				}
			}
		}
		values[i] = v
	}

	return values
}

// Result the result returned from an intermediary request
// This will return true if all of the expectations were met
func (f *Fact) Result() bool {
	if len(f.Attestations) < 1 {
		return false
	}

	for _, a := range f.Attestations {
		jws, err := jose.ParseSigned(string(a))
		if err != nil {
			return false
		}

		if !gjson.GetBytes(jws.UnsafePayloadWithoutVerification(), f.Fact).Bool() {
			return false
		}
	}

	return true
}

func (f *Fact) validate() error {
	if f.Fact == "" {
		return ErrFactEmptyName
	}

	// Skip validation if is a custom fact
	if f.Issuers != nil && len(f.Issuers) > 0 {
		return nil
	}

	for _, s := range f.Sources {
		// Return if s is not a valid source
		if _, ok := spec[s]; !ok {
			return ErrFactInvalidSource
		}

		// return error if the fact does not belong to the source
		if !contains(spec[s], f.Fact) {
			return ErrFactInvalidFactToSource
		}
	}

	if !f.hasValidOperator() {
		return ErrFactInvalidOperator
	}

	return nil
}

func (f *Fact) hasValidOperator() bool {
	var validOperators = []string{"", OperatorEqual, OperatorDifferent, OperatorGreaterOrEqualThan, OperatorGreaterThan, OperatorLessOrEqualThan, OperatorLessThan}

	for _, b := range validOperators {
		if b == f.Operator {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
