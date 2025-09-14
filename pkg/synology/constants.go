package synology

const (
	// DSM 7.2+ Container Manager constants
	ServiceName    = "pkg-ContainerManager-dockerd"
	ConfigPath     = "/var/packages/ContainerManager/etc/dockerd.json"
	DockerBinary   = "/usr/local/bin/docker"
	SocketPath     = "/var/run/docker.sock"
	RestartCommand = "sudo systemctl restart pkg-ContainerManager-dockerd"
	DefaultVolume  = "/volume1/docker"

	// SSH constants
	DefaultSSHPort = 22
	DefaultSSHUser = "admin"

	// Docker defaults
	DefaultRestartPolicy = "unless-stopped"
	DefaultNetwork       = "bridge"
)
