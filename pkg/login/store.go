package login

import "os"

func GetStoreFilePath(basedir string) (string, error) {
	if basedir[len(basedir)-1:] != "/" {
		basedir = basedir + "/"
	}
	dir := basedir + ".easyexec"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir + "/auths", nil
}
