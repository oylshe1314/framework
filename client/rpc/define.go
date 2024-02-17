package rpc

type MultiResult[result any] struct {
	Res result
	Err error
}

type MultiResults[result any] map[uint32]*MultiResult[result]
