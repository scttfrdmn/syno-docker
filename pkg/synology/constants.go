package synology

const (
	// ServiceName is the Container Manager service name in DSM 7.2+
	ServiceName = "pkg-ContainerManager-dockerd"
	// ConfigPath is the path to Container Manager configuration
	ConfigPath = "/var/packages/ContainerManager/etc/dockerd.json"
	// DockerBinary is the full path to the Docker binary
	DockerBinary = "/usr/local/bin/docker"
	// SocketPath is the path to the Docker socket
	SocketPath = "/var/run/docker.sock"
	// RestartCommand is the command to restart Container Manager
	RestartCommand = "sudo systemctl restart pkg-ContainerManager-dockerd"
	// DefaultVolume is the default volume path
	DefaultVolume = "/volume1/docker"

	// DefaultSSHPort is the default SSH port for Synology
	DefaultSSHPort = 22
	// DefaultSSHUser is the default SSH username
	DefaultSSHUser = "admin"

	// DefaultRestartPolicy is the default container restart policy
	DefaultRestartPolicy = "unless-stopped"
	// DefaultNetwork is the default Docker network mode
	DefaultNetwork = "bridge"
)
