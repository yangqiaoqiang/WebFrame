package gin

import (
	"strings"
)

type node struct {
	pattern  string  //是否为完整的url,否则为空
	part     string  //前缀树结点value
	children []*node //该节点下子节点
	isWild   bool    //实现动态路由的判断规则，:xxx 或*xxx为模糊匹配
}

func (n *node) String() string {
	return ""
}

//matchChild用于结点insert时判断此节点是否存在。存在返回该节点，否则返回nil
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		//for range遍历子节点
		if child.part == part || child.isWild {
			//找到子节点或模糊匹配 返回该节点
			return child
		}
	}
	return nil
}

//insert插入前缀树中没有的结点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		//终止条件，此时pattern为完整url
		n.pattern = pattern
		//记录在结点
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		//未匹配上，假如children子节点
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
	//插入下一节点
}

//与matchChild差别在于 matchChildren返回所有匹配的子节点,用于search。
func (n *node) matchChildren(part string) []*node {
	nodes := []*node{}
	for _, child := range n.children {
		if child.part == part || child.isWild {
			//匹配到添加入nodes
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//查找所有匹配的路径
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//递归终止条件,末尾或者通配符*
		if n.pattern == "" {
			//pattern不是完整的url，匹配失败
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		//每条路径都接着用下一part查找
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

//将所有完整的url存入list中
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
