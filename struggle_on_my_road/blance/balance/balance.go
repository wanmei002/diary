package balance

type Balancer interface {
	DoBalance([]*Instance) (*Instance, error)
}
