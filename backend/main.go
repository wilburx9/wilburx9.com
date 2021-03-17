package main

import "github.com/wilburt/wilburx9.dev/backend/common"

func main() {
	common.SetUpLogger(false)
	defer common.Logger.Sync()
}
