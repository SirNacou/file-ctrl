package app

type Handlers struct {
	GetStatus             *GetStatusHandler
	Login                 *LoginHandler
	Logout                *LogoutHandler
	OIDCLogin             *OIDCLoginHandler
	RemoveExpiredSessions *RemoveExpiredSessionsHandler
	SetupAdmin            *SetupAdminHandler
	VerifySession         *VerifySessionHandler
}
