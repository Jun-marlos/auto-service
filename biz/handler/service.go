// Code generated by hertz generator.

package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudwego/goapi/biz/dal"
	"github.com/cloudwego/goapi/biz/encrypt"
	httpclient "github.com/cloudwego/goapi/biz/http-client"
	"github.com/cloudwego/goapi/biz/mail"
	util "github.com/cloudwego/goapi/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
)

// Ping .
func Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, utils.H{
		"message": "pong",
	})
}
func Login(ctx context.Context, c *app.RequestContext) {
	IP := c.RemoteAddr().String()
	IP = util.DeletePort(IP)
	email := c.PostForm("email")
	pwd := c.PostForm("pwd")
	var uname, sessionId string
	var error_code int
	error_code = SUCCESS
	uname, sessionId, error_code = doRealLogin(ctx, email, pwd, IP)
	if sessionId != "" {
		c.SetCookie("session_id", sessionId, 24*60*60, "/", "z-coding.cn", protocol.CookieSameSiteLaxMode, false, false)
	}
	c.JSON(200, utils.H{
		"uname":      uname,
		"error_code": error_code,
	})
	if uname != "" {
		c.Redirect(304, []byte("http://txcloud.z-coding.cn/homepage"))
	} else {
		c.Redirect(304, []byte("http://txcloud.z-coding.cn/autolife"))
	}
	_, uid := dal.QueryUserInfo(email)
	dal.WriteLog(uid, int64(error_code), "用户于IP "+IP+" 登录")
}
func Logout(ctx context.Context, c *app.RequestContext) {
	error_code := SUCCESS
	defer func() {
		c.SetCookie("session_id", "", 0, "/", "z-coding.cn", protocol.CookieSameSiteLaxMode, false, false)
		c.JSON(200, utils.H{
			"error_code": error_code,
		})
		c.Redirect(304, []byte("http://txcloud.z-coding.cn/homepage"))
	}()

	uid := ctx.Value("uid")
	err := dal.RedisDelete(ctx, "[uid2sessionId]"+strconv.FormatInt(uid.(int64), 10))
	if err != nil {
		error_code = REDIS_ERROR
	}
}
func UserRegister(ctx context.Context, c *app.RequestContext) {
	token := c.PostForm("token")
	verify_code := c.PostForm("verify_code")
	uname := c.PostForm("uname")
	pwd := c.PostForm("pwd")
	error_code := SUCCESS
	if token == "" {
		token = string(c.Cookie("token"))
	}
	if token == "" || verify_code == "" || uname == "" || pwd == "" {
		error_code = BAD_PARAM
	} else if !dal.CheckVerifyCode(ctx, token, verify_code) {
		error_code = VERIFY_ERROR
	} else {
		error_code = dal.AddNewUser(ctx, token, verify_code, uname, pwd)
	}
	c.JSON(200, utils.H{
		"error_code": error_code,
	})
}
func ChangePwd(ctx context.Context, c *app.RequestContext) {

}

func AhrRegister(ctx context.Context, c *app.RequestContext) {
	student_id := c.PostForm("student_id")
	student_pwd := c.PostForm("student_pwd")
	error_code := SUCCESS
	uid := ctx.Value("uid")
	defer func() {
		dal.WriteLog(uid.(int64), int64(error_code), "注册与"+student_id+"关联的学号关系")
		c.JSON(200, utils.H{
			"error_code": error_code,
		})
	}()
	if student_id == "" || student_pwd == "" {
		error_code = BAD_PARAM
		return
	}
	err := httpclient.PassportValid(student_id, student_pwd)
	if err != nil {
		error_code = EXPECT_FALSE
		return
	}
	tmp_pwd, errr := encrypt.EncrptAESHEX([]byte(student_pwd), []byte(dal.GetConfig("encrept_key")))
	if errr != nil {
		error_code = SERVER_ERROR
		return
	}
	err = dal.AddAhrUser(ctx, student_id, string(tmp_pwd), uid.(int64))
	if err != nil {
		error_code = SERVER_ERROR
		return
	}
}
func EmailVerify(ctx context.Context, c *app.RequestContext) {
	email := c.PostForm("email")
	error_code := SUCCESS
	token := ""
	var err error
	defer func() {
		c.SetCookie("token", token, 5*60, "/", "z-coding.cn", protocol.CookieSameSiteLaxMode, false, false)
		c.JSON(200, utils.H{
			"token":      token,
			"error_code": error_code,
		})
	}()
	if email == "" {
		error_code = BAD_PARAM
		return
	}
	token, err = mail.SendVerifyCode(ctx, email)
	if err != nil {
		token = ""
		error_code = EMAIL_SEND_ERROR
		return
	}
}

// 如果session有效返回session对应的uid
func SessionValid(ctx context.Context, session_id, IP string) (int64, error) {
	uidkey := "[sessionId2uid]" + session_id + IP
	uid, err := dal.RedisQuery(ctx, uidkey)
	fmt.Println(uidkey, uid)
	if err != nil {
		return 0, err
	}
	sessionkey := "[uid2sessionId]" + uid
	realsession, err := dal.RedisQuery(ctx, sessionkey)
	fmt.Println(sessionkey)
	if err != nil || session_id+IP != realsession {
		return 0, errors.New("session_id and uid not match")
	}
	return strconv.ParseInt(uid, 10, 64)
}
func SessionVerify(ctx context.Context, c *app.RequestContext) {
	IP := c.RemoteAddr().String()
	IP = util.DeletePort(IP)
	sid := string(c.Cookie("session_id"))
	fail_url := c.Query("url")
	success_url := c.Query("success")
	fmt.Println(success_url)
	error_code := SUCCESS
	defer func() {
		if success_url != "" {
			c.Redirect(304, []byte("http://txcloud.z-coding.cn/"+success_url))
		}
	}()
	if len(sid) == 0 {
		error_code = USER_NEED_LOGIN
		c.JSON(200, utils.H{
			"error_code": error_code,
		})
		if fail_url != "" {
			c.Redirect(304, []byte("http://txcloud.z-coding.cn/"+fail_url))
		}
		return
	}
	uid, err := SessionValid(ctx, sid, IP)
	if err != nil {
		c.SetCookie("session_id", "", 0, "/", "z-coding.cn", protocol.CookieSameSiteLaxMode, false, false)
		error_code = SESSION_NOT_VALID
		c.JSON(200, utils.H{
			"error_code": error_code,
		})
		if fail_url != "" {
			c.Redirect(304, []byte("http://txcloud.z-coding.cn/"+fail_url))
		}
		return
	}
	uname := dal.QueryUserName(uid)
	c.JSON(200, utils.H{
		"uname":      uname,
		"error_code": error_code,
	})
}

func StuidQuery(ctx context.Context, c *app.RequestContext) {
	uid := ctx.Value("uid")
	sid := dal.QueryUidBindStuid(uid.(int64))
	error_code := SUCCESS
	if len(sid) == 0 {
		error_code = EXPECT_FALSE
	}
	c.JSON(200, utils.H{
		"stuid":      sid,
		"error_code": error_code,
	})
}

func SessionVerifyMiddleware(ctx context.Context, c *app.RequestContext) (context.Context, error) {
	IP := c.RemoteAddr().String()
	IP = util.DeletePort(IP)
	sid := string(c.Cookie("session_id"))
	error_code := SUCCESS
	if len(sid) == 0 {
		error_code = USER_NEED_LOGIN
		c.JSON(200, utils.H{
			"error_code": error_code,
		})
		c.Redirect(304, []byte("http://txcloud.z-coding.cn/autolife"))
		return ctx, errors.New("need login")
	}
	uid, err := SessionValid(ctx, sid, IP)
	if err != nil {
		c.SetCookie("session_id", "", 0, "/", "z-coding.cn", protocol.CookieSameSiteLaxMode, false, false)
		error_code = SESSION_NOT_VALID
		c.JSON(200, utils.H{
			"error_code": error_code,
		})
		c.Redirect(304, []byte("http://txcloud.z-coding.cn/autolife"))
		return ctx, errors.New("session invalid")
	}
	ctx = context.WithValue(ctx, "uid", uid)
	return ctx, nil
}

func doRealLogin(ctx context.Context, email, pwd, IP string) (string, string, int) {
	if email == "" || pwd == "" {
		return "", "", BAD_PARAM
	}
	if pwd != dal.QueryUserPwd(email) {
		return "", "", PWD_ERROR
	}
	uname, uid := dal.QueryUserInfo(email)
	sessionId, error_code := makeSessionId(ctx, uid, IP)
	return uname, sessionId, error_code
}

func makeSessionId(ctx context.Context, uid int64, IP string) (string, int) {
	sessionId := util.RandomStringCreate()
	value := sessionId + IP
	key1 := "[uid2sessionId]" + strconv.FormatInt(uid, 10)
	err1 := dal.RedisAdd(ctx, key1, value, time.Hour*12)
	key2 := "[sessionId2uid]" + value
	err2 := dal.RedisAdd(ctx, key2, strconv.FormatInt(uid, 10), time.Hour*12)
	if err1 != nil || err2 != nil {
		return "", REDIS_ERROR
	}
	return sessionId, SUCCESS
}

func ReportOnce(ctx context.Context, c *app.RequestContext) {
	uid := ctx.Value("uid")
	sid := dal.QueryUidBindStuid(uid.(int64))
	error_code := SUCCESS
	count := 0
	if len(sid) == 0 {
		error_code = EXPECT_FALSE
	}
	for _, one := range sid {
		err := autoReportByStuid(one)
		if err == nil {
			count++
		}
	}
	c.JSON(200, utils.H{
		"success_count": count,
		"error_code":    error_code,
	})
	dal.WriteLog(uid.(int64), int64(error_code), "完成了"+strconv.FormatInt(int64(count), 10)+"个账号的打卡")
}

func RemoveBind(ctx context.Context, c *app.RequestContext) {
	uid := ctx.Value("uid")
	err := dal.DeleteAhrUser(uid.(int64))
	error_code := SUCCESS
	if err != nil {
		error_code = SERVER_ERROR
	}
	c.JSON(200, utils.H{
		"error_code": error_code,
	})
	dal.WriteLog(uid.(int64), int64(error_code), "移除了绑定学号关系")
}

func autoReportByStuid(sid string) error {
	ahr := dal.GetReportInfo(sid)
	realpwd, err := encrypt.DecrptAESHEX(ahr.StudentPwd, []byte(dal.GetConfig("encrept_key")))
	if err != nil {
		return err
	}
	err = httpclient.AutoReport(ahr.StudentId, realpwd)
	error_code := SUCCESS
	if err != nil {
		error_code = EXPECT_FALSE
	}
	dal.WriteLog(dal.QueryUidByStuid(sid), int64(error_code), "完成对账号"+sid+"的打卡")
	return err
}

func GetLogs(ctx context.Context, c *app.RequestContext) {
	uid := ctx.Value("uid")
	js := dal.QueryLogByUid(uid.(int64))
	c.JSON(200, utils.H{
		"ops": js,
	})
}
