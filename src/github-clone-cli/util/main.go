package util

import "os"

func GetCLIUser() string {
	//TODO: Improve this
	return os.Getenv("CLI_USER_NAME")
}

func GetCLITable() string {
	//TODO: Improve this
	return os.Getenv("GITHUB_TABLE_NAME")
}
