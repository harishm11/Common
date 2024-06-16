package models

type RateRoutinSteps struct {
	RoutineId          string
	CoverageCode       string
	StepNo             int
	RefStepNo          int
	RatingItemCode     string
	RatingItemGrpCode  string
	StepOperation      string
	StepCalcMethod     string
	StepSplMethod      string
	CvgCodetoGetFctr   string
	RateVar1Code       string
	RateVar2Code       string
	RateVar3Code       string
	RateVar4Code       string
	RateVar5Code       string
	RateVar6Code       string
	RateVar7Code       string
	RateVar8Code       string
	DefaultValue       float32
	RoundorTrunc       string
	RoundorTruncDigits int
}
