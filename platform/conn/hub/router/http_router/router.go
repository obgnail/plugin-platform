package http_router

type RouterInfo struct {
	Type         string
	Method       string
	FunctionName string
	InstanceUUID string
}

type Router struct {
	// type、method、node
	trees map[string]map[string]*Node
}

func New() *Router {
	return &Router{trees: make(map[string]map[string]*Node)}
}

func (router *Router) RangeRoute(f func(Type, method string, route *RouterInfo)) {
	for Type := range router.trees {
		for method, node := range router.trees[Type] {
			node.DeepRange(func(route *RouterInfo) {
				f(Type, method, route)
			})
		}
	}
}

func (router *Router) AddRoute(Type, method, url string, handle *RouterInfo) error {
	m1 := router.trees[Type]
	if m1 == nil {
		router.trees[Type] = make(map[string]*Node)
	}
	m2 := m1[method]
	if m2 == nil {
		router.trees[Type][method] = &Node{}
	}

	return router.trees[Type][method].addRoute(url, handle)
}

func (router *Router) GetRouter(Type, method, url string) *RouterInfo {
	m1 := router.trees[Type]
	if m1 == nil {
		return nil
	}
	m2 := m1[method]
	if m2 == nil {
		return nil
	}

	handle, param, isMatch := m2.GetValue(url)

	if isMatch && handle != nil {
		if param != nil {
			releaseParam(param)
		}
		return handle
	}
	return nil
}

func (router *Router) DeepCopy() *Router {
	newMap := make(map[string]map[string]*Node, len(router.trees))

	for key1 := range router.trees {
		newMap[key1] = make(map[string]*Node)
		for key2, val2 := range router.trees[key1] {
			newMap[key1][key2] = val2.DeepCopy()
		}
	}

	return &Router{trees: newMap}
}
