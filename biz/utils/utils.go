package utils

import (
	"math/rand"
	"strings"
	"time"
)

func DeletePort(IP string) string {
	IPslice := strings.Split(IP, ":")
	return IPslice[0]
}

func RandomStringCreate() string {
	FULLSTRING := "bcdefghjklmpqrstvwyz0123456789"
	var rs []byte
	rand.Seed(time.Now().Unix())
	for i := 0; i < 30; i++ {
		rs = append(rs, FULLSTRING[rand.Int()%len(FULLSTRING)])
	}
	return (string)(rs)
}

func RandomCodeCreate() string {
	FULLSTRING := "0123456789"
	var rs []byte
	rand.Seed(time.Now().Unix())
	for i := 0; i < 6; i++ {
		rs = append(rs, FULLSTRING[rand.Int()%len(FULLSTRING)])
	}
	return (string)(rs)
}
