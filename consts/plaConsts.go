package consts

import "mitkid_web/model"
const (
	FREE_TRIAL_PLAN = 1
	THREE_MONTHS_PLAN = 2
	SIX_MONTHS_PLAN = 3
	NINE_MONTHS_PLAN = 4
	TWELVE_MONTHS_PLAN = 5
)

var PlanMap = map[int]model.Plan{
	FREE_TRIAL_PLAN:    model.Plan{PlanId: FREE_TRIAL_PLAN, PlanName: "试听课套餐", PlanTotalClass: 1, PlanPrice: 0, PlanValidity: 96},
	THREE_MONTHS_PLAN:  model.Plan{PlanId: THREE_MONTHS_PLAN, PlanName: "3个月套餐(国际双师课)", PlanTotalClass: 24, PlanPrice: 4288, PlanValidity: 5},
	SIX_MONTHS_PLAN:    model.Plan{PlanId: SIX_MONTHS_PLAN, PlanName: "6个月套餐(国际双师课)", PlanTotalClass: 48, PlanPrice: 9288, PlanValidity: 9},
	NINE_MONTHS_PLAN:   model.Plan{PlanId: NINE_MONTHS_PLAN, PlanName: "9个月套餐(国际双师课)", PlanTotalClass: 72, PlanPrice: 12588, PlanValidity: 12},
	TWELVE_MONTHS_PLAN: model.Plan{PlanId: TWELVE_MONTHS_PLAN, PlanName: "12个月套餐(国际双师课)", PlanTotalClass: 96, PlanPrice: 17288, PlanValidity: 15},
}