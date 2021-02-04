package handler

import (
	"context"
	"taobao/jackbao/user/domain/model"
	"taobao/jackbao/user/domain/service"
	user "taobao/jackbao/user/proto/user"
)

type User struct {
	UserDataService service.IUserDataService
}

//注册
func (u *User) Register(ctx context.Context, UserRegisterRequest *user.UserRegisterRequest, UserRegisterResponse *user.UserRegisterResponse) error {
	userRegister := &model.User{
		UserName:     UserRegisterRequest.UserName,
		FirstName:    UserRegisterRequest.FirstName,
		HashPassword: UserRegisterRequest.Password,
	}
	_, err := u.UserDataService.AddUser(userRegister)
	if err != nil {
		return err
	}
	UserRegisterResponse.Message = "添加成功"
	return nil
}

//登录
func (u *User) Login(ctx context.Context, UserLoginRequest *user.UserLoginRequest, UserLoginResponse *user.UserLoginResponse) error {
	isOk, err := u.UserDataService.CheckPwd(UserLoginRequest.UserName, UserLoginRequest.Password)
	if err != nil {
		return err
	}
	UserLoginResponse.IsSuccess = isOk
	return nil
}

//查询用户信息
func (u *User) GetUserInfo(ctx context.Context, UserInfoRequest *user.UserInfoRequest, UserInfoResponse *user.UserInfoResponse) error {
	userInfo, err := u.UserDataService.FindUserByName(UserInfoRequest.UserName)
	if err != nil {
		return err
	}
	UserInfoResponse = UserForResponse(userInfo)
	return nil
}

//类型转化
func UserForResponse(userModel *model.User) *user.UserInfoResponse {
	response := &user.UserInfoResponse{}
	response.UserName = userModel.UserName
	response.FirstName = userModel.FirstName
	response.UserId = userModel.ID
	return response
}
