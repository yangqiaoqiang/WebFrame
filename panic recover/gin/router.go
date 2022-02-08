package gin

import (
	"net/http"
	"strings"
)

//将Trie应用到路由中
type router struct {
	roots map[string]*node
	//添加roots存储前缀树,根节点value设为请求方式
	handlers map[string]HandlerFunc
	//router存储规则不变
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

//分割pattern字符串
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := []string{}
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	//将以method为根节点，parts为子节点insert
	r.handlers[key] = handler
}
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	//找出传入url的子节点
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			//index:结点value ,part:结点指针
			//首字母与: *匹配
			//将动态路由与实际url匹配并存入params
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//获得全部完整url的最后子节点
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := []*node{}
	root.travel(&nodes)
	return nodes
}
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		//fmt.Println(n.pattern)
		//fmt.Println(c.Path)
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
