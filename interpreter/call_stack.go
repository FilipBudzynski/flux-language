package interpreter

type CallStack struct {
	elem map[string]int
}

func (cs *CallStack) Push(funcName string) {
	if count, exists := cs.elem[funcName]; exists {
		cs.elem[funcName] = count + 1
	} else {
		cs.elem[funcName] = 1
	}
}

func (cs *CallStack) Pop(funcName string) {
	if count, exists := cs.elem[funcName]; exists {
		if count > 1 {
			cs.elem[funcName] = count - 1
		} else {
			delete(cs.elem, funcName)
		}
	}
}

func (cs *CallStack) RecursionDepth(funcName string) int {
	if count, exists := cs.elem[funcName]; exists {
		return count
	}
	return 0
}
