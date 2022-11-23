package dal

import (
	"context"
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

func QueryUserInfo(email string) (string, int64) {
	var uname string
	var uid int64
	db.Table("users").Select("uname").Where("email = ?", email).Find(&uname)
	db.Table("users").Select("uid").Where("email = ?", email).Find(&uid)
	return uname, uid
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

func AddNewUser(ctx context.Context, token, code, uname, pwd string) bool {
	token = "[token2email]" + token
	email, err := RedisQuery(ctx, token)
	if err != nil {
		return false
	}
	newuser := User{
		Uname: uname,
		Pwd:   pwd,
		Email: email,
	}
	result := db.Create(&newuser)
	if result.Error != nil {
		return false
	}
	return true
}
