package zinc_sdk

type CondExpress struct {
	logicOpValue LogicOp
	not          bool
	esCondList   []EsQueryCondition
	expressList  []*CondExpress
}

func NewAndCondExpress(not bool) *CondExpress {
	return &CondExpress{
		logicOpValue: LogicOpAnd,
		not:          not,
	}
}
func NewOrCondExpress(not bool) *CondExpress {
	return &CondExpress{
		logicOpValue: LogicOpOr,
		not:          not,
	}
}
func (ce *CondExpress) AddCond(cond EsQueryCondition) {
	if cond == nil {
		return
	}
	ce.esCondList = append(ce.esCondList, cond)
}

func (ce *CondExpress) AddExpress(ex *CondExpress) {
	if ex == nil {
		return
	}
	ce.expressList = append(ce.expressList, ex)
}
func (ce *CondExpress) ToEsQueryCondition() EsQueryCondition {
	size := len(ce.esCondList) + len(ce.expressList)
	if size == 0 {
		return nil
	}
	if size == 1 {
		var last EsQueryCondition
		if len(ce.esCondList) == 1 {
			last = ce.esCondList[0]
		} else {
			last = ce.expressList[0].ToEsQueryCondition()
		}
		if ce.not {
			return buildNotCond(last)
		}
		return last
	}

	esCondList := make([]EsQueryCondition, 0, size)
	for _, cond := range ce.esCondList {
		esCondList = append(esCondList, cond)
	}
	for _, express := range ce.expressList {
		cond := express.ToEsQueryCondition()
		if cond == nil {
			continue
		}
		esCondList = append(esCondList, cond)
	}

	if len(esCondList) == 0 {
		return nil
	}
	if len(esCondList) == 1 {
		if ce.not {
			return buildNotCond(esCondList[0])
		}
		return esCondList[0]
	}

	boolCond := &BoolCondition{}
	if ce.logicOpValue == LogicOpAnd {
		if ce.not {
			boolCond.AddMustNot(esCondList...)
		} else {
			boolCond.AddMust(esCondList...)
		}
		return boolCond
	}

	boolCond.AddShould(esCondList...)
	boolCond.BoolCond.MinShouldMatch = 1
	if ce.not {
		return buildNotCond(boolCond)
	}
	return boolCond
}
