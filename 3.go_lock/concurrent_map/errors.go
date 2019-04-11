package concurrent_map

import "fmt"

//代表参数非法的错误
type IllegalParameterError struct {
	msg string
}

func newIllegalParameterError(errMsg string) IllegalParameterError {
	return IllegalParameterError{
		msg: fmt.Sprintf("concurrent map: illegal parameter: %s", errMsg),
	}
}

func (n IllegalParameterError) Error() string {
	return n.msg
}

type IllegalPairTypeError struct {
	msg string
}

func newIllegalPairTypeError(pair Pair) IllegalPairTypeError {
	return IllegalPairTypeError{
		msg: fmt.Sprintf("concurrent map: illegal pair type: %T", pair),
	}
}

func (n IllegalPairTypeError) Error() string {
	return n.msg
}

type PairRedistributorError struct {
	msg string
}

func newPairRedistributorError(errMsg string) PairRedistributorError {
	return PairRedistributorError{
		msg: fmt.Sprintf("concurrent map: failing pair redistribution: %s", errMsg),
	}
}

func (n PairRedistributorError) Error() string {
	return n.msg
}
