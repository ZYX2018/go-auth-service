package arg

type SaveOrUpdateUserArg struct {
	Id string `json:"id"`

	Username string `json:"username"`

	Password string `json:"password"`

	UserOrigin int `json:"userOrigin"` //0：自行注册 1：邀请注册 2： 管理员手动添加

	RegisterMethod int `json:"registerMethod"` // 0：邮箱注册 1 ：手机注册

	InviteCode string `json:"inviteCode"`

	ClientId string `json:"clientId"`
}
