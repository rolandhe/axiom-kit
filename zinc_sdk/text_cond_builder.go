package zinc_sdk

type TextTypeCondBuilder struct {
	eqCondBuilder[string]
}

func (tcb *TextTypeCondBuilder) CanRange() bool {
	return false
}
func (tcb *TextTypeCondBuilder) CanMatch() bool {
	return true
}
func (tcb *TextTypeCondBuilder) CanLike() bool {
	return true
}

func (tcb *TextTypeCondBuilder) CanPrefix() bool {
	return true
}

func (tcb *TextTypeCondBuilder) NotIn(fieldName string, values ...string) EsQueryCondition {
	return buildNotCond(tcb.In(fieldName, values...))
}
func (tcb *TextTypeCondBuilder) Like(fieldName string, value string) EsQueryCondition {
	if len(value) == 0 {
		return nil
	}
	return &MatchLikeCondition{
		MatchLike: map[string]any{
			fieldName: value,
		},
	}
}
func (tcb *TextTypeCondBuilder) NotLike(fieldName string, value string) EsQueryCondition {
	return buildNotCond(tcb.Like(fieldName, value))
}
func (tcb *TextTypeCondBuilder) Match(fieldName string, value string) EsQueryCondition {
	if len(value) == 0 {
		return nil
	}
	return &MatchCondition{
		Match: map[string]any{
			fieldName: value,
		},
	}
}
func (tcb *TextTypeCondBuilder) NotMatch(fieldName string, value string) EsQueryCondition {
	return buildNotCond(tcb.Match(fieldName, value))
}
func (tcb *TextTypeCondBuilder) Prefix(fieldName string, value string) EsQueryCondition {
	if len(value) == 0 {
		return nil
	}
	return &PrefixCondition{
		Prefix: map[string]any{
			fieldName: value,
		},
	}
}
func (tcb *TextTypeCondBuilder) NotPrefix(fieldName string, value string) EsQueryCondition {
	return buildNotCond(tcb.Prefix(fieldName, value))
}

func (tcb *TextTypeCondBuilder) Gt(fieldName string, value string) EsQueryCondition {
	return nil
}
func (tcb *TextTypeCondBuilder) Gte(fieldName string, value string) EsQueryCondition {
	return nil
}
func (tcb *TextTypeCondBuilder) Lt(fieldName string, value string) EsQueryCondition {
	return nil
}
func (tcb *TextTypeCondBuilder) Lte(fieldName string, value string) EsQueryCondition {
	return nil
}
func (tcb *TextTypeCondBuilder) Between(fieldName string, lower, higher string) EsQueryCondition {
	return nil
}
