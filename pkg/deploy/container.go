package deploy

import (
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"

	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

// ContainerOptions defines options for container deployment
type ContainerOptions struct {
	Image       string
	Name        string
	Ports       []string // ["8080:80", "443:443"]
	Volumes     []string // ["/volume1/data:/app/data"]
	Env         []string // ["KEY=value"]
	Restart     string   // "unless-stopped"
	NetworkMode string
	WorkingDir  string
	Command     []string
	User        string
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	Status string
	Ports  []string
}

// NewContainerOptions creates new container options with defaults
func NewContainerOptions(image string) *ContainerOptions {
	return &ContainerOptions{
		Image:       image,
		Restart:     synology.DefaultRestartPolicy,
		NetworkMode: synology.DefaultNetwork,
	}
}

// Container deploys a container using direct Docker commands over SSH
func Container(conn *synology.Connection, opts *ContainerOptions) (string, error) {
	if opts.Name == "" {
		opts.Name = generateContainerName(opts.Image)
	}

	// Build docker run command
	dockerArgs := []string{"run", "-d"}

	// Add name
	dockerArgs = append(dockerArgs, "--name", opts.Name)

	// Add ports
	for _, port := range opts.Ports {
		dockerArgs = append(dockerArgs, "-p", port)
	}

	// Add volumes
	for _, volume := range opts.Volumes {
		dockerArgs = append(dockerArgs, "-v", volume)
	}

	// Add environment variables
	for _, env := range opts.Env {
		dockerArgs = append(dockerArgs, "-e", env)
	}

	// Add restart policy
	if opts.Restart != "" {
		dockerArgs = append(dockerArgs, "--restart", opts.Restart)
	}

	// Add network
	if opts.NetworkMode != "" {
		dockerArgs = append(dockerArgs, "--network", opts.NetworkMode)
	}

	// Add user
	if opts.User != "" {
		dockerArgs = append(dockerArgs, "--user", opts.User)
	}

	// Add working directory
	if opts.WorkingDir != "" {
		dockerArgs = append(dockerArgs, "-w", opts.WorkingDir)
	}

	// Add image
	dockerArgs = append(dockerArgs, opts.Image)

	// Add command
	dockerArgs = append(dockerArgs, opts.Command...)

	// Pull image first
	fmt.Printf("Pulling image %s...\n", opts.Image)
	if _, err := conn.ExecuteDockerCommand([]string{"pull", opts.Image}); err != nil {
		return "", errors.Wrap(err, "failed to pull image")
	}

	// Run container
	fmt.Printf("Creating and starting container %s...\n", opts.Name)
	output, err := conn.ExecuteDockerCommand(dockerArgs)
	if err != nil {
		return "", errors.Wrapf(err, "failed to run container: %s", output)
	}

	containerID := strings.TrimSpace(output)
	return containerID, nil
}

// ListContainers lists containers using direct Docker commands
func ListContainers(conn *synology.Connection, all bool) ([]ContainerInfo, error) {
	args := []string{"ps", "--format", "'table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'"}
	if all {
		args = append(args, "-a")
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list containers")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 {
		return []ContainerInfo{}, nil
	}

	var containers []ContainerInfo
	for _, line := range lines[1:] { // Skip header
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			container := ContainerInfo{
				ID:     fields[0],
				Name:   fields[1],
				Image:  fields[2],
				Status: fields[3],
			}
			if len(fields) > 4 {
				container.Ports = fields[4:]
			}
			containers = append(containers, container)
		}
	}

	return containers, nil
}

// RemoveContainer removes a container using direct Docker commands
func RemoveContainer(conn *synology.Connection, nameOrID string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, nameOrID)

	fmt.Printf("Removing container %s...\n", nameOrID)
	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to remove container: %s", output)
	}

	return nil
}

func generateContainerName(image string) string {
	// Extract image name without registry and tag
	parts := strings.Split(image, "/")
	name := parts[len(parts)-1]

	// Remove tag if present
	if idx := strings.Index(name, ":"); idx != -1 {
		name = name[:idx]
	}

	return name
}

// GetDockerClient returns a placeholder - not used in simple implementation
func GetDockerClient() (*client.Client, error) {
	// This is a placeholder for compatibility
	// The simple implementation uses SSH commands directly
	return nil, fmt.Errorf("docker client not needed in simple implementation")
}

// TestDockerConnection tests Docker availability over SSH
func TestDockerConnection(conn *synology.Connection) error {
	// Test Docker command
	if _, err := conn.ExecuteDockerCommand([]string{"version", "--format", "'{{.Server.Version}}'"}); err != nil {
		return fmt.Errorf("docker connection test failed: %w", err)
	}

	return nil
}

// ExecOptions defines options for executing commands in containers
type ExecOptions struct {
	Interactive bool
	TTY         bool
	User        string
	WorkingDir  string
	Env         []string
}

// GetContainerLogs retrieves container logs
func GetContainerLogs(conn *synology.Connection, nameOrID, tail, since string, timestamps bool) (string, error) {
	args := []string{"logs"}

	if tail != "all" && tail != "" {
		args = append(args, "--tail", tail)
	}
	if since != "" {
		args = append(args, "--since", since)
	}
	if timestamps {
		args = append(args, "--timestamps")
	}

	args = append(args, nameOrID)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get logs for container %s", nameOrID)
	}

	return output, nil
}

// FollowContainerLogs follows container logs in real-time
func FollowContainerLogs(conn *synology.Connection, nameOrID, tail, since string, timestamps bool, stdout, stderr io.Writer) error {
	args := []string{"logs", "--follow"}

	if tail != "all" && tail != "" {
		args = append(args, "--tail", tail)
	}
	if since != "" {
		args = append(args, "--since", since)
	}
	if timestamps {
		args = append(args, "--timestamps")
	}

	args = append(args, nameOrID)

	cmd := strings.Join(append([]string{"/usr/local/bin/docker"}, args...), " ")
	return conn.StreamCommand(cmd, stdout, stderr)
}

// ExecCommand executes a command in a container and returns output
func ExecCommand(conn *synology.Connection, nameOrID string, command []string, opts *ExecOptions) (string, error) {
	args := []string{"exec"}

	if opts.User != "" {
		args = append(args, "--user", opts.User)
	}
	if opts.WorkingDir != "" {
		args = append(args, "--workdir", opts.WorkingDir)
	}
	for _, env := range opts.Env {
		args = append(args, "--env", env)
	}

	args = append(args, nameOrID)
	args = append(args, command...)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute command in container %s", nameOrID)
	}

	return output, nil
}

// ExecInteractive executes an interactive command in a container
func ExecInteractive(conn *synology.Connection, nameOrID string, command []string, opts *ExecOptions) error {
	args := []string{"exec"}

	if opts.Interactive {
		args = append(args, "-i")
	}
	if opts.TTY {
		args = append(args, "-t")
	}
	if opts.User != "" {
		args = append(args, "--user", opts.User)
	}
	if opts.WorkingDir != "" {
		args = append(args, "--workdir", opts.WorkingDir)
	}
	for _, env := range opts.Env {
		args = append(args, "--env", env)
	}

	args = append(args, nameOrID)
	args = append(args, command...)

	cmd := strings.Join(append([]string{"/usr/local/bin/docker"}, args...), " ")
	return conn.StreamCommand(cmd, nil, nil)
}

// RestartContainer restarts a container
func RestartContainer(conn *synology.Connection, nameOrID string, timeout int) error {
	args := []string{"restart"}
	if timeout > 0 {
		args = append(args, "--time", fmt.Sprintf("%d", timeout))
	}
	args = append(args, nameOrID)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to restart container %s: %s", nameOrID, output)
	}

	return nil
}

// StartContainer starts a stopped container
func StartContainer(conn *synology.Connection, nameOrID string) error {
	args := []string{"start", nameOrID}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to start container %s: %s", nameOrID, output)
	}

	return nil
}

// StopContainer stops a running container
func StopContainer(conn *synology.Connection, nameOrID string, timeout int) error {
	args := []string{"stop"}
	if timeout > 0 {
		args = append(args, "--time", fmt.Sprintf("%d", timeout))
	}
	args = append(args, nameOrID)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to stop container %s: %s", nameOrID, output)
	}

	return nil
}

// StatsOptions defines options for container stats
type StatsOptions struct {
	All      bool
	NoStream bool
	Format   string
}

// ImageInfo represents Docker image information
type ImageInfo struct {
	Repository string
	Tag        string
	Digest     string
	ID         string
	Created    string
	Size       string
}

// ImagesOptions defines options for listing images
type ImagesOptions struct {
	All       bool
	Dangling  bool
	Digests   bool
	Format    string
	NoTrunc   bool
	Quiet     bool
}

// PullOptions defines options for pulling images
type PullOptions struct {
	AllTags             bool
	Platform            string
	Quiet               bool
	DisableContentTrust bool
}

// RmiOptions defines options for removing images
type RmiOptions struct {
	Force   bool
	NoPrune bool
}

// SystemDfItem represents disk usage information
type SystemDfItem struct {
	Type        string
	Total       string
	Active      string
	Size        string
	Reclaimable string
}

// SystemDfOptions defines options for system df
type SystemDfOptions struct {
	Format  string
	Verbose bool
}

// SystemInfoOptions defines options for system info
type SystemInfoOptions struct {
	Format string
}

// SystemPruneOptions defines options for system prune
type SystemPruneOptions struct {
	All     bool
	Force   bool
	Volumes bool
	Filter  []string
}

// SystemPruneResult represents the result of system prune
type SystemPruneResult struct {
	ContainersDeleted int
	ImagesDeleted     int
	NetworksDeleted   int
	VolumesDeleted    int
	SpaceReclaimed    string
}

// ShowContainerStats displays container resource usage statistics
func ShowContainerStats(conn *synology.Connection, containers []string, opts *StatsOptions) error {
	args := []string{"stats"}

	if opts.All {
		args = append(args, "--all")
	}
	if opts.NoStream {
		args = append(args, "--no-stream")
	}
	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	}

	args = append(args, containers...)

	cmd := strings.Join(append([]string{"/usr/local/bin/docker"}, args...), " ")
	return conn.StreamCommand(cmd, nil, nil)
}

// ListImages lists Docker images
func ListImages(conn *synology.Connection, repository string, opts *ImagesOptions) ([]ImageInfo, error) {
	args := []string{"images"}

	if opts.All {
		args = append(args, "--all")
	}
	if opts.Dangling {
		args = append(args, "--filter", "dangling=true")
	}
	if opts.Digests {
		args = append(args, "--digests")
	}
	if opts.NoTrunc {
		args = append(args, "--no-trunc")
	}

	if opts.Digests {
		args = append(args, "--format", "table {{.Repository}}\t{{.Tag}}\t{{.Digest}}\t{{.ID}}\t{{.CreatedSince}}\t{{.Size}}")
	} else {
		args = append(args, "--format", "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedSince}}\t{{.Size}}")
	}

	if repository != "" {
		args = append(args, repository)
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list images")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 {
		return []ImageInfo{}, nil
	}

	var images []ImageInfo
	for _, line := range lines[1:] { // Skip header
		fields := strings.Fields(line)
		if opts.Digests && len(fields) >= 6 {
			image := ImageInfo{
				Repository: fields[0],
				Tag:        fields[1],
				Digest:     fields[2],
				ID:         fields[3],
				Created:    fields[4],
				Size:       fields[5],
			}
			images = append(images, image)
		} else if !opts.Digests && len(fields) >= 5 {
			image := ImageInfo{
				Repository: fields[0],
				Tag:        fields[1],
				ID:         fields[2],
				Created:    fields[3],
				Size:       fields[4],
			}
			images = append(images, image)
		}
	}

	return images, nil
}

// ListImageIDs lists Docker image IDs only
func ListImageIDs(conn *synology.Connection, repository string, opts *ImagesOptions) ([]string, error) {
	args := []string{"images", "--quiet"}

	if opts.All {
		args = append(args, "--all")
	}
	if opts.Dangling {
		args = append(args, "--filter", "dangling=true")
	}
	if opts.NoTrunc {
		args = append(args, "--no-trunc")
	}

	if repository != "" {
		args = append(args, repository)
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list image IDs")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var imageIDs []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			imageIDs = append(imageIDs, line)
		}
	}

	return imageIDs, nil
}

// PullImage pulls a Docker image
func PullImage(conn *synology.Connection, imageName string, opts *PullOptions) error {
	args := []string{"pull"}

	if opts.AllTags {
		args = append(args, "--all-tags")
	}
	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}
	if opts.Quiet {
		args = append(args, "--quiet")
	}
	if opts.DisableContentTrust {
		args = append(args, "--disable-content-trust")
	}

	args = append(args, imageName)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to pull image %s: %s", imageName, output)
	}

	if !opts.Quiet {
		fmt.Print(output)
	}

	return nil
}

// RemoveImage removes a Docker image
func RemoveImage(conn *synology.Connection, imageName string, opts *RmiOptions) error {
	args := []string{"rmi"}

	if opts.Force {
		args = append(args, "--force")
	}
	if opts.NoPrune {
		args = append(args, "--no-prune")
	}

	args = append(args, imageName)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to remove image %s: %s", imageName, output)
	}

	return nil
}

// GetSystemDf gets Docker system disk usage
func GetSystemDf(conn *synology.Connection, opts *SystemDfOptions) ([]SystemDfItem, error) {
	args := []string{"system", "df"}

	if opts.Verbose {
		args = append(args, "--verbose")
	}
	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	} else {
		args = append(args, "--format", "table {{.Type}}\t{{.Total}}\t{{.Active}}\t{{.Size}}\t{{.Reclaimable}}")
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get system disk usage")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 {
		return []SystemDfItem{}, nil
	}

	var usage []SystemDfItem
	for _, line := range lines[1:] { // Skip header
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			item := SystemDfItem{
				Type:        fields[0],
				Total:       fields[1],
				Active:      fields[2],
				Size:        fields[3],
				Reclaimable: fields[4],
			}
			usage = append(usage, item)
		}
	}

	return usage, nil
}

// GetSystemInfo gets Docker system information
func GetSystemInfo(conn *synology.Connection, opts *SystemInfoOptions) (string, error) {
	args := []string{"system", "info"}

	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return "", errors.Wrap(err, "failed to get system info")
	}

	return output, nil
}

// SystemPrune removes unused Docker data
func SystemPrune(conn *synology.Connection, opts *SystemPruneOptions) (*SystemPruneResult, error) {
	args := []string{"system", "prune"}

	if opts.All {
		args = append(args, "--all")
	}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.Volumes {
		args = append(args, "--volumes")
	}
	for _, filter := range opts.Filter {
		args = append(args, "--filter", filter)
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prune system")
	}

	// Parse the output to extract numbers (simplified parsing)
	result := &SystemPruneResult{
		SpaceReclaimed: "unknown",
	}

	// This is a simplified implementation - in practice you'd parse the actual output
	// to extract specific numbers for containers, images, networks, volumes deleted
	fmt.Print(output)

	return result, nil
}

// VolumeInfo represents Docker volume information
type VolumeInfo struct {
	Name   string
	Driver string
}

// VolumeListOptions defines options for listing volumes
type VolumeListOptions struct {
	Format string
	Quiet  bool
}

// VolumeCreateOptions defines options for creating volumes
type VolumeCreateOptions struct {
	Driver  string
	Labels  []string
	Options []string
}

// VolumeRemoveOptions defines options for removing volumes
type VolumeRemoveOptions struct {
	Force bool
}

// VolumeInspectOptions defines options for inspecting volumes
type VolumeInspectOptions struct {
	Format string
}

// VolumePruneOptions defines options for pruning volumes
type VolumePruneOptions struct {
	Force  bool
	Filter []string
}

// VolumePruneResult represents the result of volume prune
type VolumePruneResult struct {
	VolumesDeleted int
	SpaceReclaimed string
}

// InspectOptions defines options for inspecting objects
type InspectOptions struct {
	Format string
	Size   bool
	Type   string
}

// ExportOptions defines options for exporting containers
type ExportOptions struct {
	Output string
}

// ImportOptions defines options for importing images
type ImportOptions struct {
	Change   []string
	Message  string
	Platform string
}

// ListVolumes lists Docker volumes
func ListVolumes(conn *synology.Connection, opts *VolumeListOptions) ([]VolumeInfo, error) {
	args := []string{"volume", "ls"}

	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	} else {
		args = append(args, "--format", "table {{.Driver}}\t{{.Name}}")
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list volumes")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 {
		return []VolumeInfo{}, nil
	}

	var volumes []VolumeInfo
	for _, line := range lines[1:] { // Skip header
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			volume := VolumeInfo{
				Driver: fields[0],
				Name:   fields[1],
			}
			volumes = append(volumes, volume)
		}
	}

	return volumes, nil
}

// ListVolumeNames lists Docker volume names only
func ListVolumeNames(conn *synology.Connection, opts *VolumeListOptions) ([]string, error) {
	args := []string{"volume", "ls", "--quiet"}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list volume names")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var volumeNames []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			volumeNames = append(volumeNames, line)
		}
	}

	return volumeNames, nil
}

// CreateVolume creates a Docker volume
func CreateVolume(conn *synology.Connection, volumeName string, opts *VolumeCreateOptions) (string, error) {
	args := []string{"volume", "create"}

	if opts.Driver != "" {
		args = append(args, "--driver", opts.Driver)
	}
	for _, label := range opts.Labels {
		args = append(args, "--label", label)
	}
	for _, option := range opts.Options {
		args = append(args, "--opt", option)
	}

	if volumeName != "" {
		args = append(args, volumeName)
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create volume: %s", output)
	}

	return strings.TrimSpace(output), nil
}

// RemoveVolume removes a Docker volume
func RemoveVolume(conn *synology.Connection, volumeName string, opts *VolumeRemoveOptions) error {
	args := []string{"volume", "rm"}

	if opts.Force {
		args = append(args, "--force")
	}

	args = append(args, volumeName)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to remove volume %s: %s", volumeName, output)
	}

	return nil
}

// InspectVolume inspects a Docker volume
func InspectVolume(conn *synology.Connection, volumeName string, opts *VolumeInspectOptions) (string, error) {
	args := []string{"volume", "inspect"}

	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	}

	args = append(args, volumeName)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return "", errors.Wrapf(err, "failed to inspect volume %s", volumeName)
	}

	return output, nil
}

// PruneVolumes removes unused Docker volumes
func PruneVolumes(conn *synology.Connection, opts *VolumePruneOptions) (*VolumePruneResult, error) {
	args := []string{"volume", "prune"}

	if opts.Force {
		args = append(args, "--force")
	}
	for _, filter := range opts.Filter {
		args = append(args, "--filter", filter)
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prune volumes")
	}

	// Parse the output to extract numbers (simplified parsing)
	result := &VolumePruneResult{
		SpaceReclaimed: "unknown",
	}

	// This is a simplified implementation - in practice you'd parse the actual output
	fmt.Print(output)

	return result, nil
}

// InspectObject inspects a Docker object (container, image, volume, network)
func InspectObject(conn *synology.Connection, objectName string, opts *InspectOptions) (string, error) {
	args := []string{"inspect"}

	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	}
	if opts.Size {
		args = append(args, "--size")
	}
	if opts.Type != "" {
		args = append(args, "--type", opts.Type)
	}

	args = append(args, objectName)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return "", errors.Wrapf(err, "failed to inspect object %s", objectName)
	}

	return output, nil
}

// ExportContainer exports a container's filesystem as a tar archive
func ExportContainer(conn *synology.Connection, containerName string, opts *ExportOptions) error {
	args := []string{"export"}

	if opts.Output != "" {
		args = append(args, "--output", opts.Output)
	}

	args = append(args, containerName)

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to export container %s: %s", containerName, output)
	}

	if opts.Output == "" {
		fmt.Print(output)
	}

	return nil
}

// ImportImage imports the contents from a tarball to create a filesystem image
func ImportImage(conn *synology.Connection, source, repository string, opts *ImportOptions) (string, error) {
	args := []string{"import"}

	for _, change := range opts.Change {
		args = append(args, "--change", change)
	}
	if opts.Message != "" {
		args = append(args, "--message", opts.Message)
	}
	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}

	args = append(args, source)
	if repository != "" {
		args = append(args, repository)
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return "", errors.Wrapf(err, "failed to import image from %s: %s", source, output)
	}

	return strings.TrimSpace(output), nil
}
