package login

func GetStoreFilePath(basedir string) string {
	if basedir[len(basedir)-1:] != "/" {
		basedir = basedir + "/"
	}
	dir := basedir + "Easyexec"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir + "/auths", nil
}
