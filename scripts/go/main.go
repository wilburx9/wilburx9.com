package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	scripts "github.com/wilburt/wilburx9.dev/scripts/tools"
)

func main() {
	pass := uuid.NewString()
	salt := uuid.NewString()
	var data = map[string]string{
		"pass": pass,
		"salt": salt,
		"hash": scripts.GenerateHash(pass, salt),
	}
	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(string(j))
}
