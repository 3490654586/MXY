package MXY

import "strings"

type Node struct {
	pattern string //待匹配的路由
	part string //路由中的一部分
	children []*Node //子节点
	isWild bool
}

// 第一个匹配成功的节点，用于插入
func (n *Node) matchChild(part string) *Node {
	//循环子节点
	for _, child := range n.children {
		//如果当前结点的的child参数和传进来的参数相同child或者为精准匹配说明里面有个次路由,则返回此结点
		if child.part == part || child.isWild {
			return child
		}
	}
	//没有返回nil
	return nil
}
// 所有匹配成功的节点，用于查找
func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	//循环子节点
	for _, child := range n.children {
		//当前字节的路由=part传进来的路由或者为静态匹配,添加到nodes里面返回
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入
func (n *Node) insert(pattern string, parts []string, height int) {
	//如果parts的长度 == 树的高度,说明没有待匹配的路由返回
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	//
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		//新增一个结点
		child = &Node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		//将新增结点添加到当前结点
		n.children = append(n.children, child)
	}
	//递归插入
	child.insert(pattern, parts, height+1)
}

//查询
func (n *Node) search(parts []string, height int) *Node {
	//如果len的长度==高度说明次次路由没有下一个结点,或者路由前面带*通配符
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		//递归查询,height+1代表查询下一个结点
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}