package service

import (
	"mitkid_web/consts"
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
			planMap[planItem.AccountId] = append(planMap[planItem.AccountId], planItem)
		} else {
			listc = make([]model.AccountPlan, 0)
			listc = append(listc, planItem)
			planMap[planItem.AccountId] = listc
		}
	}
	return
}
func FullPlan(a *model.AccountPlan) {
	plan := planConsts.PLAN_MAP[a.PlanCode]
	a.PlanName = plan.PlanName
	a.PlanTotalClass = plan.PlanTotalClass
}
func FullPlanList(list []model.AccountPlan) {
	for i, _ := range list {
		a := &list[i]
		plan := planConsts.PLAN_MAP[a.PlanCode]
		a.PlanName = plan.PlanName
		a.PlanTotalClass = plan.PlanTotalClass
	}

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
	planCreatedAt := time.Now()
	ap := &model.AccountPlan{
		AccountId:     id,
		PlanCode:      p.PlanCode,
		PlanCreatedAt: planCreatedAt,
		Status:        consts.PLAN_NOACTIVE_STATUS,
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
	aps, err = s.dao.ListPlanByPlanIds(pIds)
	FullPlanList(aps)
	return aps, err
}
func (s *Service) BatchUpdatePlanUsedClass(accountId string, planMap map[int]int) (err error) {
	for k, v := range planMap {
		if err = s.dao.BatchUpdatePlanUsedClass(accountId, k, v); err != nil {
			return err
		}
	}
	return nil
}
