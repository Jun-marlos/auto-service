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

func RedisAdd(ctx context.Context, key, value string) int {
	err := rdb.Set(ctx, key, value, time.Hour*24).Err()
	if err != nil {
		return 301
	}
	return 0
}
