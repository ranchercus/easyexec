package common

var (
	RANCHER_URL = "https://rancher.i.fbank.com/v3-public/localProviders/local?action=login"
	RANCHER_kUBECONFIG_URL = "https://10.2.10.29:31443/v3/clusters/local?action=generateKubeconfig"
	//RANCHER_kUBECONFIG_URL = "https://rancher.i.fbank.com/v3/clusters/local?action=generateKubeconfig"
	RANCHER_NAMESPACES        = "https://rancher.i.fbank.com/v3/cluster/local/namespaces?limit=-1&sort=name"
)
