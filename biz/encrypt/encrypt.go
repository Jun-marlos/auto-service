package encrypt

import (
	"fmt"
	"os"

	"github.com/robertkrimen/otto"
)

func EncryptWithSalt(raw, salt string) string {
	jsfile := "/usr/local/go/src/github.com/cloudwego/goapi/biz/encrypt/encrypt.js"
	bytes, err := os.ReadFile(jsfile)
	if err != nil {
		return "err1"
	}
	vm := otto.New()
	_, err = vm.Run(bytes)
	if err != nil {
		return "err2"
	}
	enc, errr := vm.Call("encryptAES", nil, raw, salt)
	fmt.Println(enc)
	if errr != nil {
		return ""
	}
	return enc.String()
}
