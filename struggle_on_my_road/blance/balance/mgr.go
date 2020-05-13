package balance

import "errors"

type BalanceMgr struct {
	allBalancer map[string]Balancer
}

var Mgr = BalanceMgr{
	allBalancer: make(map[string]Balancer),
}

func (blc *BalanceMgr) RegisterMgr(name string, b Balancer) {
	blc.allBalancer[name] = b
}

func RegisterMgr(name string, b Balancer){
	Mgr.RegisterMgr(name, b)
}

func DoBalance(name string, insts []*Instance) (inst *Instance, err error) {
	balancer , ok := Mgr.allBalancer[name]
	if !ok {
		err = errors.New("not found "+ name +" balance")
		return
	}
	inst, err = balancer.DoBalance(insts)
	return
}