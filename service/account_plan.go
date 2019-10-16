package service

import (
	"mitkid_web/consts"
	"mitkid_web/model"
)

var PlanMap = map[int]model.Plan{
	consts.FREE_TRIAL_PLAN:    model.Plan{PlanCode: consts.FREE_TRIAL_PLAN, PlanName: "试听课套餐", PlanTotalClass: 1, PlanPrice: 0, PlanValidity: 96},
	consts.THREE_MONTHS_PLAN:  model.Plan{PlanCode: consts.THREE_MONTHS_PLAN, PlanName: "3个月套餐(国际双师课)", PlanTotalClass: 24, PlanPrice: 4288, PlanValidity: 5},
	consts.SIX_MONTHS_PLAN:    model.Plan{PlanCode: consts.SIX_MONTHS_PLAN, PlanName: "6个月套餐(国际双师课)", PlanTotalClass: 48, PlanPrice: 9288, PlanValidity: 9},
	consts.NINE_MONTHS_PLAN:   model.Plan{PlanCode: consts.NINE_MONTHS_PLAN, PlanName: "9个月套餐(国际双师课)", PlanTotalClass: 72, PlanPrice: 12588, PlanValidity: 12},
	consts.TWELVE_MONTHS_PLAN: model.Plan{PlanCode: consts.TWELVE_MONTHS_PLAN, PlanName: "12个月套餐(国际双师课)", PlanTotalClass: 96, PlanPrice: 17288, PlanValidity: 15},
}

func (s *Service) ListAccountPlansWithAccountIDs(accountIds []string) (planMap map[string]([]model.AccountPlan), err error) {
	plans, err := s.dao.ListAccountPlansWithAccountIDs(accountIds)
	if err != nil {
		return nil, err
	}
	if plans == nil {
		return nil, nil
	}
	planMap = make(map[string]([]model.AccountPlan))
	for _, planItem := range plans {
		FullPlan(&planItem)
		if listc, ok := planMap[planItem.AccountId]; ok {
			listc = append(listc, planItem)
		} else {
			listc = make([]model.AccountPlan, 0)
			listc = append(listc, planItem)
			planMap[planItem.AccountId] = listc
		}
	}
	return
}
func FullPlan(a *model.AccountPlan) {
	plan := PlanMap[a.PlanCode]
	a.PlanName = plan.PlanName
	a.PlanTotalClass = plan.PlanTotalClass
}
