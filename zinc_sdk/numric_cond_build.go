package zinc_sdk

type NumRicCondBuild[T RangeType] struct {
	noTextCondBuilder[T]
	eqCondBuilder[T]
}

func (n *NumRicCondBuild[T]) Gt(fieldName string, value T) EsQueryCondition {
	rg := &RangeCondition[T]{}
	rgField := &RangeField[T]{
		Gt: &value,
	}
	rg.AddCondition(fieldName, rgField)
	return rg
}
func (n *NumRicCondBuild[T]) Gte(fieldName string, value T) EsQueryCondition {
	rg := &RangeCondition[T]{}
	rgField := &RangeField[T]{
		Gte: &value,
	}
	rg.AddCondition(fieldName, rgField)
	return rg
}
func (n *NumRicCondBuild[T]) Lt(fieldName string, value T) EsQueryCondition {
	rg := &RangeCondition[T]{}
	rgField := &RangeField[T]{
		Lt: &value,
	}
	rg.AddCondition(fieldName, rgField)
	return rg
}
func (n *NumRicCondBuild[T]) Lte(fieldName string, value T) EsQueryCondition {
	rg := &RangeCondition[T]{}
	rgField := &RangeField[T]{
		Lte: &value,
	}
	rg.AddCondition(fieldName, rgField)
	return rg
}
func (n *NumRicCondBuild[T]) Between(fieldName string, lower, higher T) EsQueryCondition {
	rg := &RangeCondition[T]{}
	rgField := &RangeField[T]{
		Gte: &lower,
		Lte: &higher,
	}
	rg.AddCondition(fieldName, rgField)
	return rg
}
