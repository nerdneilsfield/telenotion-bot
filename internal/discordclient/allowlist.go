package discordclient

func IsAllowedUser(allowed []string, userID string) bool {
	if userID == "" {
		return false
	}
	for _, id := range allowed {
		if id == userID {
			return true
		}
	}
	return false
}
