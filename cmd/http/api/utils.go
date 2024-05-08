package api

func uuidPath(path string) string {
	var bracketDepth int
	for i := 0; i < len(path); i++ {
		if path[i] == '{' {
			bracketDepth++
		}
		if path[i] == '}' {
			bracketDepth--
		}
		if bracketDepth == 1 && path[i] == ':' {
			if path[i+1:i+5] == "uuid" {
				path = path[:i] + ":[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}" + path[i+5:] //nolint:lll
			}
		}
	}
	return path
}
