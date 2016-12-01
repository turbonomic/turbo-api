package api

// targetType, nameOrAddress, targetIdentifier, username, password
type Target struct {
	TargetType       string
	NameOrAddress    string
	TargetIdentifier string

	Username string
	Password string
}
