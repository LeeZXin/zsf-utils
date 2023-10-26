package selector

import (
	"errors"
)

type MasterSlavesSelector[T any] struct {
	master        Node[T]
	slaves        []Node[T]
	slaveSelector Selector[T]
}

func NewMasterSlavesSelector[T any](lbPolicy string, nodes ...Node[T]) (*MasterSlavesSelector[T], error) {
	selectorFunc, b := FindNewSelectorFunc[T](lbPolicy)
	if !b {
		selectorFunc = NewRoundRobinSelector[T]
	}
	if len(nodes) == 0 {
		return nil, errors.New("empty nodes")
	}
	master := nodes[0]
	var slaves []Node[T]
	if len(nodes) > 1 {
		slaves = nodes[1:]
	}
	var ret Selector[T]
	if len(slaves) == 0 {
		ret = newErrorSelector[T](EmptyNodesErr)
	} else {
		ret = selectorFunc(slaves)
	}
	return &MasterSlavesSelector[T]{
		master:        master,
		slaves:        slaves,
		slaveSelector: ret,
	}, nil
}

func (m *MasterSlavesSelector[T]) Master() T {
	return m.master.Data
}

func (m *MasterSlavesSelector[T]) Slaves() []T {
	ret := make([]T, 0, len(m.slaves))
	for _, slave := range m.slaves {
		ret = append(ret, slave.Data)
	}
	return ret
}

func (m *MasterSlavesSelector[T]) SelectSlave(keys ...string) (T, error) {
	node, err := m.slaveSelector.Select(keys...)
	if err != nil {
		var ret T
		return ret, err
	}
	return node.Data, nil
}

func (m *MasterSlavesSelector[T]) IndexSlave(index int) (T, error) {
	if index >= len(m.slaves) {
		var ret T
		return ret, errors.New("index out of bound")
	}
	return m.slaves[index].Data, nil
}
