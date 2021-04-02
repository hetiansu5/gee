package gee

import (
	"strings"
)

type node struct {
	pattern  string         // 待匹配路由，例如 /p/:lang
	part     string         // 路由中的一部分，例如 :lang
	children []*node        // 子节点，例如 [doc, tutorial, intro]
	isWild   bool           // 是否精确匹配，part 含有 : 或 * 时为true
	params   map[int]string //路由参数
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		for index, part := range parts {
			if isWild(part) && len(part[1:]) > 0 {
				if len(n.params) == 0 {
					n.params = map[int]string{}
				}
				n.params[index] = part[1:]
			}
		}
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{pattern: "", part: part, isWild: isWild(part)}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	nodes := n.matchChildren(part)

	for _, nod := range nodes {
		result := nod.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) parseParams(parts []string) map[string]string {
	if len(n.params) == 0 {
		return nil
	}

	l := len(parts)
	params := map[string]string{}
	for index, key := range n.params {
		if l > index {
			params[key] = parts[index]
		}
	}
	return params
}

func isWild(part string) bool {
	return part[0] == '*' || part[0] == ':'
}

func parsePattern(pattern string) []string {
	arr := strings.Split(pattern, "/")
	var parts []string
	for _, part := range arr {
		if part == "" {
			continue
		}
		parts = append(parts, part)
	}
	return parts
}
