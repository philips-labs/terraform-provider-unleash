package provider

func getUserRoleId(role string) int {
	rolesMap := make(map[string]int)
	rolesMap["Admin"] = 1
	rolesMap["Editor"] = 2
	rolesMap["Viewer"] = 3
	return rolesMap[role]
}

func getUserRole(roleId int) string {
	rolesMap := make(map[int]string)
	rolesMap[1] = "Admin"
	rolesMap[2] = "Editor"
	rolesMap[3] = "Viewer"
	return rolesMap[roleId]
}
