package graph

type Graph map[string][]string

func (g Graph) BFS(from, to string) ([]string, bool) {
	if from == to {
		return []string{from}, true
	}

	visited := map[string]bool{from: true}
	parent := map[string]string{}
	queue := []string{from}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, next := range g[cur] {
			if visited[next] {
				continue
			}
			visited[next] = true
			parent[next] = cur
			if next == to {
				return reconstruct(parent, from, to), true
			}
			queue = append(queue, next)
		}
	}

	return nil, false
}

func (g Graph) DFS(from, to string) ([]string, bool) {
	visited := map[string]bool{}
	parent := map[string]string{}
	var found bool

	var walk func(string)
	walk = func(node string) {
		if found {
			return
		}
		visited[node] = true
		if node == to {
			found = true
			return
		}
		for _, next := range g[node] {
			if visited[next] {
				continue
			}
			parent[next] = node
			walk(next)
			if found {
				return
			}
		}
	}

	walk(from)
	if !found {
		return nil, false
	}
	return reconstruct(parent, from, to), true
}

func reconstruct(parent map[string]string, from, to string) []string {
	path := []string{to}
	for cur := to; cur != from; {
		cur = parent[cur]
		if cur == "" {
			return nil
		}
		path = append([]string{cur}, path...)
	}
	return path
}
