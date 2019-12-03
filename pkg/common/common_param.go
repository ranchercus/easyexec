package common

import "os"

type CommonType struct {
	PodName        string
	DeploymentName string
	Namespace      string
	ContainerIndex int
}

func getBaseDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = os.TempDir()
	}
	if dir[len(dir)-1:] != "/" {
		dir = dir + "/"
	}
	return dir
}

func GetConfigPath() (string, error) {
	dir := getBaseDir()
	path, err := GetDir(dir)
	if err != nil {
		return "", err
	}
	path = path + "/config"
	return path, nil
}

func GetCookiePath() (string, error) {
	dir := getBaseDir()
	path, err := GetDir(dir)
	if err != nil {
		return "", err
	}
	path = path + "/cookie"
	return path, nil
}
func GetDir(basedir string) (string, error) {
	dir := basedir + Config_Dir
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir, nil
}