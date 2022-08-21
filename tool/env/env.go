package env

import "os"

const (
	KEY_DEPLOY  = "deploy"
	VAL_DEV     = "dev"
	VAL_PREVIEW = "preview"
	VAL_PRODUCT = "product"
)

func IsProduct() bool {
	deploy, ok := os.LookupEnv(KEY_DEPLOY)
	if !ok {
		return false
	}
	return deploy == VAL_PRODUCT
}
