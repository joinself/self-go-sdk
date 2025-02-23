package message

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo darwin LDFLAGS: -lself_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#cgo linux LDFLAGS: -lself_sdk
#include <self-sdk.h>
#include <stdlib.h>
*/
import "C"

const (
	OperatorEquals              ComparisonOperator = C.OPERATOR_EQUALS
	OperatorNotEquals           ComparisonOperator = C.OPERATOR_NOT_EQUALS
	OperatorGreaterThan         ComparisonOperator = C.OPERATOR_GREATER_THAN
	OperatorGreaterThanOrEquals ComparisonOperator = C.OPERATOR_GREATER_THAN_OR_EQUALS
	OperatorLessThan            ComparisonOperator = C.OPERATOR_LESS_THAN
	OperatorLessThanOrEquals    ComparisonOperator = C.OPERATOR_LESS_THAN_OR_EQUALS
)

type ComparisonOperator int

func (c ComparisonOperator) String() string {
	switch c {
	case OperatorEquals:
		return "Equals"
	case OperatorNotEquals:
		return "Not Equals"
	case OperatorGreaterThan:
		return "Greater Than"
	case OperatorGreaterThanOrEquals:
		return "Greater Than Or Equals"
	case OperatorLessThan:
		return "Less Than"
	case OperatorLessThanOrEquals:
		return "Less Than Or Equals"
	default:
		return "Unknown"
	}
}
