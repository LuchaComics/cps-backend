package constants

type key int

const (
	SessionIsAuthorized key = iota
	SessionSkipAuthorization
	SessionID
	SessionIPAddress
	SessionUser
	SessionUserTenantID
	SessionUserRoleID
	SessionUserID
	SessionUserUUID
	SessionUserTimezone
	SessionUserFullName
)
