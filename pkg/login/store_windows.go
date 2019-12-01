package login

import "os"

func getDir(basedir string) (string, error) {
	dir := basedir + "Easyexec"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func GetStoreFilePath(basedir string) (string, error) {
	dir, err := getDir(basedir)
	return dir + "/config", err
}

func GetStoreCookiePath(basedir string) (string, error) {
	dir, err := getDir(basedir)
	return dir + "/cookie", err
}
