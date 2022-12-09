package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	err error
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("connect redis error, error code : %v", err)
	}
	fmt.Println("connect redis successfully")

	dsn := username + ":" + pwd + "@tcp(" + host + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("connect mysql-db successfully")

}

func QueryUserPwd(email string) string {
	var pwd string
	db.Table("users").Select("pwd").Where("email = ?", email).Find(&pwd)
	return pwd
}

func QueryUserName(uid int64) string {
	var uname string
	db.Table("users").Select("uname").Where("uid = ?", uid).Find(&uname)
	return uname
}

func QueryUserEmailByUid(uid int64) string {
	var email string
	db.Table("users").Select("email").Where("uid = ?", uid).Find(&email)
	return email
}

func QueryUserInfo(email string) (string, int64) {
	var uname string
	var uid int64
	db.Table("users").Select("uname").Where("email = ?", email).Find(&uname)
	db.Table("users").Select("uid").Where("email = ?", email).Find(&uid)
	return uname, uid
}

func QueryUidBindStuid(uid int64) []string {
	var ans []string
	db.Table("ahrs").Select("student_id").Where("uid = ?", uid).Find(&ans)
	return ans
}

func QueryUidByStuid(stuid string) int64 {
	var ans int64
	db.Table("ahrs").Select("uid").Where("student_id = ?", stuid).Find(&ans)
	return ans
}

func RedisAdd(ctx context.Context, key, value string, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	return err
}

func RedisQuery(ctx context.Context, key string) (string, error) {
	result := rdb.Get(ctx, key)
	res := result.Val()
	err := result.Err()
	return res, err
}

func RedisDelete(ctx context.Context, key string) error {
	err := rdb.Del(ctx, key).Err()
	return err
}

func CheckVerifyCode(ctx context.Context, token, code string) bool {
	token = "[verify_code]" + token
	realcode, err := RedisQuery(ctx, token)
	if err == nil && code == realcode {
		RedisDelete(ctx, token)
		return true
	}
	return false
}

func AddNewUser(ctx context.Context, token, code, uname, pwd string) int {
	token = "[token2email]" + token
	email, err := RedisQuery(ctx, token)
	if err != nil {
		return 301
	}
	_, uid := QueryUserInfo(email)
	if uid != 0 {
		return 204
	}
	newuser := User{
		Uname: uname,
		Pwd:   pwd,
		Email: email,
	}
	result := db.Create(&newuser)
	if result.Error != nil {
		return 300
	}
	return 0
}

func AddAhrUser(ctx context.Context, studentid, pwd string, uid int64) error {
	newone := Ahr{}
	newone.Uid = uid
	newone.StudentPwd = pwd
	newone.StudentId = studentid
	result := db.Save(&newone)
	return result.Error
}

func GetReportInfo(stuid string) Ahr {
	rt := Ahr{}
	db.Where("student_id = ?", stuid).Find(&rt)
	return rt
}

func DeleteAhrUser(uid int64) error {
	return db.Delete(&Ahr{}, "uid = ?", uid).Error
}

func UpdateAhrUser(ctx context.Context, studentid, pwd string, uid int64) error {
	return nil
}

func GetConfig(configName string) string {
	var rt string
	db.Table("configs").Select("config_content").Where("config_name = ?", configName).Find(&rt)
	return rt
}

func WriteLog(uid, error_code int64, content string) error {
	alog := Log{
		Uid:       uid,
		ErrorCode: error_code,
		Content:   content,
	}
	return db.Create(&alog).Error
}

func QueryLogByUid(uid int64) string {
	logs := []Log{}
	db.Select("id, error_code, content, create_time").Where("uid = ?", uid).Order("id DESC").Find(&logs)
	ans, err := json.Marshal(logs)
	if err != nil {
		return ""
	}
	return string(ans)
}
