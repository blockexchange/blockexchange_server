package components

import "blockexchange/controller"

type AccessTokenModel struct {
}

func AccessToken(rc *controller.RenderContext) *AccessTokenModel {

	return &AccessTokenModel{}
}
