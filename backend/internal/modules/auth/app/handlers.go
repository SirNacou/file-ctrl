package app

type Handlers struct {
	GetStatus             *GetStatusHandler
	Login                 *LoginHandler
	Logout                *LogoutHandler
	RemoveExpiredSessions *RemoveExpiredSessionsHandler
	SetupAdmin            *SetupAdminHandler
	VerifySession         *VerifySessionHandler
}
