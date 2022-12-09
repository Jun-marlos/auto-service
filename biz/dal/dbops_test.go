package dal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

func TestCreate(t *testing.T) {
	testuser := User{
		Uname: "*****",
		Pwd:   "123456",
		Email: "*******@qq.com",
	}
	testahr := Ahr{
		Uid:        1,
		StudentId:  "******",
		StudentPwd: "******",
	}
	testahr.LastSuccessDate, _ = time.Parse(TIME_LAYOUT, "2022-11-22 16:25:23")
	testlog := Log{
		Uid:       1,
		ErrorCode: 100,
		Content:   "justfortest",
	}
	result := db.Create(&testuser)
	fmt.Println(result)
	result = db.Create(&testahr)
	fmt.Println(result)
	result = db.Create(&testlog)
	fmt.Println(result)
}

func TestQuery(t *testing.T) {
	var testuser []User
	var testahr []Ahr
	var testlog []Log
	result1 := db.Table("users").Find(&testuser)
	result2 := db.Table("ahrs").Find(&testahr)
	result3 := db.Table("logs").Find(&testlog)
	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
}

func TestRedisSet(t *testing.T) {
	ctx := context.Background()
	//Set方法的最后一个参数表示过期时间，0表示永不过期
	err = rdb.Set(ctx, "key1", "value1", 0).Err()
	if err != nil {
		panic(err)
	}

	//key2将会在两分钟后过期失效
	err = rdb.Set(ctx, "key2", "value2", time.Minute*2).Err()
	if err != nil {
		panic(err)
	}
}

func TestRedisGet(t *testing.T) {
	ctx := context.Background()
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("key: %v\n", val)

	val2, err := rdb.Get(ctx, "key-not-exist").Result()
	if err == redis.Nil {
		fmt.Println("key不存在")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("值为: %v\n", val2)
	}
}

func TestQueryPwd(t *testing.T) {
	fmt.Println(QueryUserPwd("*****@qq.com"))
	fmt.Println(QueryUserInfo("*******@qq.com"))
	fmt.Println(QueryUserPwd("*****@qq.com"))
}

func TestRedis(t *testing.T) {
	ctx := context.Background()
	RedisAdd(ctx, "just-test", "test-answer", time.Hour*1)
	ans, _ := RedisQuery(ctx, "just-test")
	fmt.Println(ans)
	RedisDelete(ctx, "just-test")
	ans, _ = RedisQuery(ctx, "just-test")
	fmt.Println(ans)
}

func TestVerifyCode(t *testing.T) {
	ctx := context.Background()
	RedisAdd(ctx, "[verify_code]abcdefg", "05488", time.Hour*1)
	fmt.Println(CheckVerifyCode(ctx, "abcdefg", "05488"))
}

func TestQueryNotExist(t *testing.T) {
	a, b := QueryUserInfo("*******@qq.com")
	fmt.Println(a, b)
}

func TestAddahr(t *testing.T) {
	ctx := context.Background()
	AddAhrUser(ctx, "2292019******", "zhaohhhhh", 8)
}

func TestConfig(t *testing.T) {
	fmt.Println(GetConfig("encrept_key"))
}

func TestDeleteAhrbyUid(t *testing.T) {
	DeleteAhrUser(8)
}

func TestGetAhrinfo(t *testing.T) {
	fmt.Println(GetReportInfo("******"))
}

func TestWriteLog(t *testing.T) {
	WriteLog(8, 11, "just for test")
}

func TestGetlog(t *testing.T) {
	fmt.Println(QueryLogByUid(8))
}
