package planConsts

import (
	"mitkid_web/consts"
	"mitkid_web/model"
)

var PLAN_MAP = map[int]model.Plan{
	consts.FREE_TRIAL_PLAN:    model.Plan{PlanCode: consts.FREE_TRIAL_PLAN, PlanName: "试听课套餐", PlanTotalClass: 1, PlanPrice: 0, PlanValidity: 96},
	consts.THREE_MONTHS_PLAN:  model.Plan{PlanCode: consts.THREE_MONTHS_PLAN, PlanName: "3个月套餐(国际双师课)", PlanTotalClass: 24, PlanPrice: 4288, PlanValidity: 5},
	consts.SIX_MONTHS_PLAN:    model.Plan{PlanCode: consts.SIX_MONTHS_PLAN, PlanName: "6个月套餐(国际双师课)", PlanTotalClass: 48, PlanPrice: 9288, PlanValidity: 9},
	consts.NINE_MONTHS_PLAN:   model.Plan{PlanCode: consts.NINE_MONTHS_PLAN, PlanName: "9个月套餐(国际双师课)", PlanTotalClass: 72, PlanPrice: 12588, PlanValidity: 12},
	consts.TWELVE_MONTHS_PLAN: model.Plan{PlanCode: consts.TWELVE_MONTHS_PLAN, PlanName: "12个月套餐(国际双师课)", PlanTotalClass: 96, PlanPrice: 17288, PlanValidity: 15},
}
