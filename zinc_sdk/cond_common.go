package zinc_sdk

import "time"

type LogicOp = int32

const (
	LogicOpAnd LogicOp = 0
	LogicOpOr  LogicOp = 1
)

type CondType interface {
	string | bool | RangeType
}

type NumericType interface {
	int64 | int32 | float64 | float32
}

type RangeType interface {
	time.Time | NumericType
}

type CondBuilder[T CondType] interface {
	CanRange() bool
	CanMatch() bool
	CanLike() bool
	CanPrefix() bool

	In(fieldName string, values ...T) EsQueryCondition
	NotIn(fieldName string, values ...T) EsQueryCondition
	Like(fieldName string, value string) EsQueryCondition
	NotLike(fieldName string, value string) EsQueryCondition
	Match(fieldName string, value string) EsQueryCondition
	NotMatch(fieldName string, value string) EsQueryCondition
	Prefix(fieldName string, value string) EsQueryCondition
	NotPrefix(fieldName string, value string) EsQueryCondition
	Eq(fieldName string, Value T) EsQueryCondition
	Neq(fieldName string, Value T) EsQueryCondition
	Gt(fieldName string, value T) EsQueryCondition
	Gte(fieldName string, value T) EsQueryCondition
	Lt(fieldName string, value T) EsQueryCondition
	Lte(fieldName string, value T) EsQueryCondition
	Between(fieldName string, lower, higher T) EsQueryCondition
}

type eqCondBuilder[T CondType] struct {
}

func (ic *eqCondBuilder[T]) Eq(fieldName string, value T) EsQueryCondition {
	return &TermCondition[string]{
		Term: map[string]any{
			fieldName: value,
		},
	}
}
func (ic *eqCondBuilder[T]) Neq(fieldName string, value T) EsQueryCondition {
	return buildNotCond(ic.Eq(fieldName, value))
}
func (ic *eqCondBuilder[T]) In(fieldName string, values ...T) EsQueryCondition {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return ic.Eq(fieldName, values[0])
	}
	return &TermsCondition[string]{
		Terms: map[string]any{
			fieldName: values,
		},
	}
}
func (ic *eqCondBuilder[T]) NotIn(fieldName string, values ...T) EsQueryCondition {
	return buildNotCond(ic.In(fieldName, values...))
}

type noTextCondBuilder[T CondType] struct {
}

func (n *noTextCondBuilder[T]) CanRange() bool {
	return true
}
func (n *noTextCondBuilder[T]) CanMatch() bool {
	return false
}
func (n *noTextCondBuilder[T]) CanLike() bool {
	return false
}

func (n *noTextCondBuilder[T]) CanPrefix() bool {
	return false
}

func (n *noTextCondBuilder[T]) Like(fieldName string, value string) EsQueryCondition {
	return nil
}
func (n *noTextCondBuilder[T]) NotLike(fieldName string, value string) EsQueryCondition {
	return nil
}
func (n *noTextCondBuilder[T]) Match(fieldName string, value string) EsQueryCondition {
	return nil
}
func (n *noTextCondBuilder[T]) NotMatch(fieldName string, value string) EsQueryCondition {
	return nil
}
func (n *noTextCondBuilder[T]) Prefix(fieldName string, value string) EsQueryCondition {
	return nil
}
func (n *noTextCondBuilder[T]) NotPrefix(fieldName string, value string) EsQueryCondition {
	return nil
}

func buildNotCond(cond EsQueryCondition) *BoolCondition {
	if cond == nil {
		return nil
	}
	boolCond := &BoolCondition{}
	boolCond.AddMustNot(cond)
	return boolCond
}
