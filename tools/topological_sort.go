package tools

import "slices"

func TopologicalSort[T comparable](edges map[T][]T) (sorted []T, cycle []T) {
	// 1. 构建入度表
	degreeMap := make(map[T]int)
	for cur, edges := range edges {
		if _, ok := degreeMap[cur]; !ok {
			degreeMap[cur] = 0
		}
		for _, edge := range edges {
			degreeMap[edge]++
		}
	}
	// 2. 构建0入度队列
	queue := make([]T, 0, len(edges))
	for name, degree := range degreeMap {
		if degree == 0 {
			queue = append(queue, name)
		}
	}
	// 3. 构建拓扑排序
	sorted = make([]T, 0, len(edges))
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		sorted = append(sorted, curr)

		for _, dep := range edges[curr] {
			degreeMap[dep]--
			if degreeMap[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	// 4. 环检测
	if len(sorted) != len(edges) {
		cycle = CheckCircular(edges)
	}

	return
}

// 检测图中的环并返回环的节点序列
func CheckCircular[T comparable](edges map[T][]T) (cycle []T) {
	visited := make(map[T]bool)
	recStack := make(map[T]bool)

	var dfs func(node T) bool
	var path []T

	dfs = func(node T) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, neighbor := range edges[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				// 发现环，提取环的信息
				cycleStartIndex := slices.Index(path, neighbor)
				if cycleStartIndex != -1 {
					cycle = append([]T(nil), path[cycleStartIndex:]...)
					cycle = append(cycle, neighbor)
				}
				return true
			}
		}

		recStack[node] = false
		path = path[:len(path)-1]
		return false
	}

	// 从所有未访问的节点开始DFS
	for node := range edges {
		if !visited[node] {
			path = []T{}
			if dfs(node) {
				break
			}
		}
	}

	return cycle
}
