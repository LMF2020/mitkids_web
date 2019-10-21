package service

import (
	"mitkid_web/consts/planConsts"
	"mitkid_web/model"
	"time"
)

func (s *Service) ListAccountPlansWithAccountIDs(accountIds []string) (planMap map[string]([]model.AccountPlan), err error) {
	plans, err := s.dao.ListAccountPlansWithAccountIDs(accountIds)
	if err != nil {
		return nil, err
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
	plan := planConsts.PlanMap[a.PlanCode]
	a.PlanName = plan.PlanName
	a.PlanTotalClass = plan.PlanTotalClass
}

func (s *Service) ListAccountPlansWithAccountID(accountId string) (plans []model.AccountPlan, err error) {
	plans, err = s.dao.ListAccountPlansWithAccountID(accountId)
	if err != nil {
		return nil, err
	}
	for i, _ := range plans {
		FullPlan(&plans[i])
	}
	return
}

func (s *Service) AddUserPlan(id string, p *model.Plan) (err error) {
	ap := &model.AccountPlan{
		AccountId:     id,
		PlanCode:      p.PlanCode,
		PlanCreatedAt: time.Now(),
		PlanExpiredAt: time.Now().AddDate(0, p.PlanValidity, 0),
	}
	return s.dao.AddUserPlan(ap)
}

func (s *Service) GetPlanByPlanId(pId int) (ap *model.AccountPlan, err error) {
	return s.dao.GetPlanByPlanId(pId)
}
func (s *Service) DeletePlanByPlanId(pId int) (err error) {
	return s.dao.DeletePlanByPlanId(pId)
}

func (s *Service) ListPlanByPlanIds(pIds []int) (aps []model.AccountPlan, err error) {
	return s.dao.ListPlanByPlanIds(pIds)
}
