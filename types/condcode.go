package types

type CondCode uint8

const (
	CondAL CondCode = iota
	CondPL
)

var CondCodesByLit = map[string]CondCode{
	"AL": CondAL,
	"PL": CondPL,
}

var CondCodeBits = map[CondCode]uint32{
	CondAL: 0b1110,
	CondPL: 0b1010,
}
