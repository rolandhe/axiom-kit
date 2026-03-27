package zinc_sdk

type BoolTypeCondBuilder struct {
	noTextCondBuilder[bool]
}

func (bcb *BoolTypeCondBuilder) CanRange() bool {
	return false
}
func (bcb *BoolTypeCondBuilder) CanMatch() bool {
	return false
}
func (bcb *BoolTypeCondBuilder) CanLike() bool {
	return false
}

func (bcb *BoolTypeCondBuilder) CanPrefix() bool {
	return false
}

func (bcb *BoolTypeCondBuilder) Eq(fieldName string, value bool) EsQueryCondition {
	return &TermCondition[string]{
		Term: map[string]any{
			fieldName: value,
		},
	}
}

func (bcb *BoolTypeCondBuilder) Neq(fieldName string, value bool) EsQueryCondition {
	return nil
}

func (bcb *BoolTypeCondBuilder) In(fieldName string, values ...bool) EsQueryCondition {
	return nil
}
func (bcb *BoolTypeCondBuilder) NotIn(fieldName string, values ...bool) EsQueryCondition {
	return nil
}
