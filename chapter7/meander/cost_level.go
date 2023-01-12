package meander

import (
	"errors"
	"strings"
)

// Cost represents a cost level.
// 表現しようとしている列挙が数種類しかないので int8 を基底型としている
type Cost int8

// Cost enumerator represents various levels of Cost.
// 列挙型 Cost を表現
const (
	_ Cost = iota // iota は const に連番の整数を割り当てる // 0 はあえて無視している
	Cost1
	Cost2
	Cost3
	Cost4
	Cost5
)

var costStrings = map[string]Cost{
	"$":     Cost1,
	"$$":    Cost2,
	"$$$":   Cost3,
	"$$$$":  Cost4,
	"$$$$$": Cost5,
}

// String Java で言う toString だーね https://pkg.go.dev/fmt#Stringer
func (l Cost) String() string {
	for s, v := range costStrings {
		if l == v {
			return s
		}
	}
	return "invalid"
}

// ParseCost parses the cost string into a Cost type.
func ParseCost(s string) Cost {
	return costStrings[s]
}

// CostRange represents a range of Cost values.
type CostRange struct {
	From Cost
	To   Cost
}

func (r CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

// ParseCostRange parses a cost range string into a CostRange.
func ParseCostRange(s string) (CostRange, error) {
	var r CostRange
	segs := strings.Split(s, "...")
	if len(segs) != 2 {
		return r, errors.New("invalid cost range")
	}
	r.From = ParseCost(segs[0])
	r.To = ParseCost(segs[1])
	return r, nil
}
