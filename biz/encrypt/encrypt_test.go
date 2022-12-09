package encrypt

import (
	"fmt"
	"testing"
)

const (
	TestPwd  = "hhhhhhhh"
	TestSalt = "urdqIjgVdxSHknGj"
	TestAns  = "mTNrcARA73dqNgbnmRet0fI6iA5s1bi+S4xxI+v0MWQN2p66m51Tng0XlAZTpUj1Kh5V/ydx5LvNYWeTnYAF1xdmnwDHnSTvC5/pU9gG3hY="
)

func TestEncrypt(t *testing.T) {
	fmt.Println(EncryptWithSalt(TestPwd, TestSalt))
}

func TestAes(t *testing.T) {
	origin := "imniuboyi"
	key := "zhao40880884zhao"
	secrect, err := EncrptAESHEX([]byte(origin), []byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(secrect)
	originde, errr := DecrptAESHEX(secrect, []byte(key))
	if errr != nil {
		fmt.Println(errr)
		return
	}
	fmt.Println(string(originde))

}
