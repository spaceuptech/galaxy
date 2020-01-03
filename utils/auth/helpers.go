package auth

func (m *Module) GetUserName() string {
	return m.config.UserName
}
