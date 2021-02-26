package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"net/http"
)

func DoSystemUserLogin(lr *LoginRequest) (*LoginResponse, int, error) {
	// get username and password
	mobile := lr.Mobile
	password := lr.Password
	clientType := lr.ClientType

	//validate username and password if is empty
	if len(mobile) == 0 || len(password) == 0 {
		return nil, http.StatusBadRequest, errors.New("手机号和密码不能为空")
	}

	// connect db
	o := orm.NewOrm()

	// check the username if existing
	user := &SystemUser{Mobile: mobile}
	err := o.Read(user, "mobile")
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("用户不存在")
	}

	// generate the password hash
	hash, err := GeneratePassHash(password, user.Salt)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if hash != user.Password {
		return nil, http.StatusBadRequest, errors.New("密码不正确")
	}

	// generate token
	tokenString, err := GenerateEncToken(lr, user.Id, 86400*10)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	filter := o.QueryTable("system_user_multi_login").Filter("userid", user.Id).Filter("client_type", lr.ClientType)

	count, _ := filter.Count()
	if count < 1 {
		o.Insert(&SystemUserMultiLogin{Token: tokenString, ClientType: clientType, UserId: user.Id})
		logs.Debug("insert token ", user.Mobile, tokenString)
	} else {
		var systemUserMultiLogin SystemUserMultiLogin
		err := filter.One(&systemUserMultiLogin)
		if err == nil {
			systemUserMultiLogin.Token = tokenString
			update, _ := o.Update(&systemUserMultiLogin)
			logs.Debug("update token  ", update, ToJson(systemUserMultiLogin))
		}
	}

	return &LoginResponse{
		Username:   user.Username,
		ClientType: lr.ClientType,
		UserID:     user.Id,
		Token:      tokenString,
	}, http.StatusOK, nil
}

// DoCreateUser: create a user
func DoCreateSystemUser(cr *CreateRequest) (*CreateResponse, int, error) {
	// connect db
	o := orm.NewOrm()

	// check username if exist
	userNameCheck := SystemUser{Mobile: cr.Mobile}
	err := o.Read(&userNameCheck, "mobile")
	if err == nil {
		return nil, http.StatusBadRequest, errors.New("手机号已经存在")
	}

	if len(cr.Password) > 16 || len(cr.Password) < 6 {
		return nil, http.StatusBadRequest, errors.New("密码长度不能小于6，不能大于16")
	}

	// generate salt
	saltKey, err := GenerateSalt()
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// generate password hash
	hash, err := GeneratePassHash(cr.Password, saltKey)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// create user
	user := SystemUser{}
	user.Username = cr.Username
	user.Mobile = cr.Mobile
	user.Password = hash
	user.Salt = saltKey

	_, err = o.Insert(&user)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, errors.New("注册失败，请稍后重试或联系客服解决")
	}

	return &CreateResponse{
		UserID:   user.Id,
		Username: user.Username,
		Mobile:   user.Mobile,
	}, http.StatusOK, nil
}

func GetSystemUser(encToken string) (*SystemUser, error) {
	jwtPayload, err := ValidateEncToken(encToken)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	userNameCheck := SystemUser{Id: jwtPayload.UserID}
	err = o.Read(&userNameCheck, "id")

	if err != nil {
		return nil, errors.New("获取用户信息失败")
	} else {
		return &SystemUser{
			Id:       userNameCheck.Id,
			Username: userNameCheck.Username,
			Mobile:   userNameCheck.Mobile,
		}, nil

	}
}
