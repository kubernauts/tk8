package provisioner

func SetClusteName() {
	if len(Name) < 1 {
		config := GetClusterConfig()
		Name = config.AwsClusterName
	}
}
