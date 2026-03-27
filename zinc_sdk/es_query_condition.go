package zinc_sdk

type EsQueryCondition = any

type ESQueryReq struct {
	Query  EsQueryCondition      `json:"query"`
	Size   int                   `json:"size"`
	Source []string              `json:"_source,omitempty"`
	Sort   []map[string]*EsOrder `json:"sort,omitempty"`
}

type EsOrder struct {
	Order string `json:"order"`
}

type BoolCondition struct {
	BoolCond struct {
		Must           []EsQueryCondition `json:"must,omitempty"`
		Should         []EsQueryCondition `json:"should,omitempty"`
		MustNot        []EsQueryCondition `json:"must_not,omitempty"`
		Filter         []EsQueryCondition `json:"filter,omitempty"`
		MinShouldMatch int                `json:"minimum_should_match,omitempty"`
	} `json:"bool"`
}

func (boolCondition *BoolCondition) AddMust(conds ...EsQueryCondition) *BoolCondition {
	for _, cond := range conds {
		boolCondition.BoolCond.Must = append(boolCondition.BoolCond.Must, cond)
	}
	return boolCondition
}

func (boolCondition *BoolCondition) AddShould(conds ...EsQueryCondition) *BoolCondition {
	for _, cond := range conds {
		boolCondition.BoolCond.Should = append(boolCondition.BoolCond.Should, cond)
	}
	return boolCondition
}

func (boolCondition *BoolCondition) AddMustNot(conds ...EsQueryCondition) *BoolCondition {
	for _, cond := range conds {
		boolCondition.BoolCond.MustNot = append(boolCondition.BoolCond.MustNot, cond)
	}
	return boolCondition
}

func (boolCondition *BoolCondition) AddFilter(cond EsQueryCondition) *BoolCondition {
	boolCondition.BoolCond.Filter = append(boolCondition.BoolCond.Filter, cond)
	return boolCondition
}

type MatchCondition struct {
	Match map[string]any `json:"match"`
}

func (match *MatchCondition) AddCondition(field string, value string) {
	if len(value) == 0 {
		return
	}
	if match.Match == nil {
		match.Match = make(map[string]any)
	}
	match.Match[field] = value
}
func (match *MatchCondition) AddConditionWithBoost(field string, value string, boost float64) {
	if boost == 0 {
		match.AddCondition(field, value)
		return
	}
	if len(value) == 0 {
		return
	}
	if match.Match == nil {
		match.Match = make(map[string]any)
	}
	valueMap := make(map[string]any)
	valueMap["boost"] = boost
	valueMap["query"] = value
	match.Match[field] = valueMap
}

type MatchPhraseCondition struct {
	MatchPhrase map[string]any `json:"match_phrase"`
}

func (mp *MatchPhraseCondition) AddCondition(field string, value string) {
	if len(value) == 0 {
		return
	}
	if mp.MatchPhrase == nil {
		mp.MatchPhrase = make(map[string]any)
	}
	mp.MatchPhrase[field] = value
}

func (mp *MatchPhraseCondition) AddConditionWithBoost(field string, value string, boost float64) {
	if boost == 0 {
		mp.AddCondition(field, value)
		return
	}
	if len(value) == 0 {
		return
	}
	if mp.MatchPhrase == nil {
		mp.MatchPhrase = make(map[string]any)
	}
	valueMap := make(map[string]any)
	valueMap["boost"] = boost
	valueMap["query"] = value
	mp.MatchPhrase[field] = valueMap
}

type MatchLikeCondition struct {
	MatchLike map[string]any `json:"match_like"`
}

func (mp *MatchLikeCondition) AddCondition(field string, value string) {
	if len(value) == 0 {
		return
	}
	if mp.MatchLike == nil {
		mp.MatchLike = make(map[string]any)
	}
	mp.MatchLike[field] = value
}

func (mp *MatchLikeCondition) AddConditionWithBoost(field string, value string, boost float64) {
	if boost == 0 {
		mp.AddCondition(field, value)
		return
	}
	if len(value) == 0 {
		return
	}
	if mp.MatchLike == nil {
		mp.MatchLike = make(map[string]any)
	}
	valueMap := make(map[string]any)
	valueMap["boost"] = boost
	valueMap["query"] = value
	mp.MatchLike[field] = valueMap
}

//type TermsValueType interface {
//	string | int64 | int | float32 | float64 | int32 | bool
//}

type TermCondition[T CondType] struct {
	Term map[string]any `json:"term"`
}

func (tc *TermCondition[T]) AddCondition(field string, value *T) {
	if value == nil {
		return
	}
	if tc.Term == nil {
		tc.Term = make(map[string]any)
	}
	tc.Term[field] = value
}

func (tc *TermCondition[T]) AddConditionWithBoost(field string, value *T, boost float64) {
	if boost == 0 {
		tc.AddCondition(field, value)
		return
	}
	if value == nil {
		return
	}
	if tc.Term == nil {
		tc.Term = make(map[string]any)
	}
	valueMap := make(map[string]any)
	valueMap["boost"] = boost
	valueMap["value"] = value
	tc.Term[field] = valueMap
}

type TermsCondition[T CondType] struct {
	Terms map[string]any `json:"terms"`
}

func (tsc *TermsCondition[T]) AddCondition(field string, value []T) {
	if len(value) == 0 {
		return
	}
	if tsc.Terms == nil {
		tsc.Terms = make(map[string]any)
	}
	tsc.Terms[field] = value
}

func (tsc *TermsCondition[T]) AddConditionWithBoost(field string, value []T, boost float64) {
	tsc.AddCondition(field, value)
	if boost == 0 {
		return
	}
	if tsc.Terms == nil {
		return
	}
	tsc.Terms["boost"] = boost
}

type PrefixCondition struct {
	Prefix map[string]any `json:"prefix"`
}

func (pc *PrefixCondition) AddCondition(field string, value string) {
	if len(value) == 0 {
		return
	}
	if pc.Prefix == nil {
		pc.Prefix = make(map[string]any)
	}
	pc.Prefix[field] = value
}
func (pc *PrefixCondition) AddConditionWithBoost(field string, value string, boost float64) {
	if boost == 0 {
		pc.AddCondition(field, value)
		return
	}
	if len(value) == 0 {
		return
	}
	if pc.Prefix == nil {
		pc.Prefix = make(map[string]any)
	}
	valueMap := make(map[string]any)
	valueMap["boost"] = boost
	valueMap["value"] = value
	pc.Prefix[field] = valueMap
}

type WildcardCondition struct {
	Wildcard map[string]any `json:"wildcard"`
}

func (wc *WildcardCondition) AddCondition(field string, value string) {
	if len(value) == 0 {
		return
	}
	if wc.Wildcard == nil {
		wc.Wildcard = make(map[string]any)
	}
	wc.Wildcard[field] = value
}
func (wc *WildcardCondition) AddConditionWithBoost(field string, value string, boost float64) {
	if boost == 0 {
		wc.AddCondition(field, value)
		return
	}
	if len(value) == 0 {
		return
	}
	if wc.Wildcard == nil {
		wc.Wildcard = make(map[string]any)
	}
	valueMap := make(map[string]any)
	valueMap["boost"] = boost
	valueMap["value"] = value
	wc.Wildcard[field] = valueMap
}

type FuzzyCondition struct {
	Fuzzy map[string]any `json:"fuzzy"`
}

func (fc *FuzzyCondition) AddCondition(field string, value string) {
	if len(value) == 0 {
		return
	}
	if fc.Fuzzy == nil {
		fc.Fuzzy = make(map[string]any)
	}
	fc.Fuzzy[field] = value
}
func (fc *FuzzyCondition) AddConditionWithBoost(field string, value string, boost float64) {
	if boost == 0 {
		fc.AddCondition(field, value)
		return
	}
	if len(value) == 0 {
		return
	}
	if fc.Fuzzy == nil {
		fc.Fuzzy = make(map[string]any)
	}
	valueMap := make(map[string]any)
	valueMap["boost"] = boost
	valueMap["value"] = value

	fc.Fuzzy[field] = valueMap
}

type RangeCondition[T RangeType] struct {
	Range map[string]any `json:"range"`
}

func (rc *RangeCondition[T]) AddCondition(field string, value *RangeField[T]) {
	if value == nil {
		return
	}
	if rc.Range == nil {
		rc.Range = map[string]any{}
	}
	rc.Range[field] = value
}

func (rc *RangeCondition[T]) AddConditionWithBoost(field string, value *RangeField[T], boost float64) {
	if boost == 0 {
		rc.AddCondition(field, value)
		return
	}
	if value == nil {
		return
	}
	if rc.Range == nil {
		rc.Range = map[string]any{}
	}
	value.Boost = boost
	rc.Range[field] = value
}

type RangeField[T RangeType] struct {
	Gt       *T      `json:"gt,omitempty"`
	Gte      *T      `json:"gte,omitempty"`
	Lt       *T      `json:"lt,omitempty"`
	Lte      *T      `json:"lte,omitempty"`
	Format   string  `json:"format,omitempty"`
	TimeZone string  `json:"time_zone,omitempty"`
	Boost    float64 `json:"boost,omitempty"`
}

//func ConvertTermValuesToAnySlice[I TermsValueType](arr []I) []any {
//	ret := make([]any, 0, len(arr))
//	for _, v := range arr {
//		ret = append(ret, v)
//	}
//	return ret
//}
