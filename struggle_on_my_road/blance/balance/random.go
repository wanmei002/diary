package balance

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {

}

func init(){
	RegisterMgr("random", &RandomBalance{})
}

func (rdm *RandomBalance) DoBalance(inst []*Instance) (insts *Instance, err error){
	if len(inst) <= 0 {
		err = errors.New("no instance")
		return
	}

	n := len(inst) - 1
	insts = inst[rand.Intn(n)]
	return
}
