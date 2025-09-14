# Remote container management on Synology DSM 7.x

The most secure and practical approach for remotely managing Docker containers on Synology DSM 7.x systems is **SSH with Docker contexts**, requiring specific PATH environment fixes and proper certificate management. DSM 7.2 introduced Container Manager (replacing the Docker package) with configuration changes at `/var/packages/ContainerManager/etc/dockerd.json`, requiring service management through systemd rather than the legacy upstart system. Three primary methods exist for remote management: SSH-based access (recommended for security), TLS-secured TCP socket connections on port 2376 (best for automation), and native Synology APIs (limited functionality).

## Critical DSM 7.x changes affect remote access

Container Manager replaced the Docker package in DSM 7.2, introducing significant architectural changes that impact remote management strategies. The configuration file relocated from `/var/packages/Docker/etc/dockerd.json` to `/var/packages/ContainerManager/etc/dockerd.json`, while the service name changed from `pkg-Docker-dockerd` to `pkg-ContainerManager-dockerd`. These changes require updated service management commands using systemd (`sudo systemctl restart pkg-ContainerManager-dockerd`) rather than the deprecated synopkgctl commands.

The Docker binary remains at `/usr/local/bin/docker` but isn't included in the default SSH PATH during non-interactive sessions, causing "docker: command not found" errors. This PATH issue affects all remote SSH executions and requires one of three workarounds: enabling user environment support in sshd_config, creating a symbolic link to /usr/bin/docker, or specifying full paths in all commands. The socket file at `/var/run/docker.sock` maintains compatibility but requires proper group permissions for non-root access.

Container Manager includes native Docker Compose support through the GUI's Project feature, though it lacks some advanced capabilities like .env file support and has performance limitations compared to third-party tools. The default shared folder `/volume1/docker/` serves as the standard location for container data, with all volume mappings requiring full path specifications rather than relative paths.

## SSH provides the most secure remote access method

SSH-based Docker management offers **optimal security** without exposing additional ports, leveraging existing authentication infrastructure and providing comprehensive audit trails. Setting up SSH access requires enabling the SSH service through Control Panel → Terminal & SNMP, preferably on a custom port between 49152-65535 to avoid automated attacks. Users must belong to the Administrators group for SSH access in DSM 7.x, a security restriction that cannot be bypassed through standard configuration.

The critical configuration step involves fixing the PATH environment issue by editing `/etc/ssh/sshd_config` to include `PermitUserEnvironment PATH`, then creating a `.ssh/environment` file containing the full PATH including `/usr/local/bin`. After applying these changes and restarting the SSH service with `sudo synoservicectl --reload sshd`, remote Docker commands function correctly. SSH key authentication provides additional security, requiring proper permissions (755 for home directory, 644 for authorized_keys) due to DSM's non-standard default permissions.

Docker contexts simplify SSH-based management through persistent configurations. Creating a context with `docker context create synology --docker host=ssh://admin@nas-ip:port` enables seamless switching between local and remote Docker environments. For enhanced performance, SSH connection multiplexing reduces overhead through ControlMaster settings in ~/.ssh/config, maintaining persistent connections for up to 10 minutes after the last command.

## TLS certificate authentication enables secure TCP connections

For automation and high-volume operations, **TLS-secured TCP connections on port 2376** provide better performance than SSH while maintaining security. Certificate generation requires creating a Certificate Authority, server certificate with proper Subject Alternative Names (including all IP addresses and hostnames), and client certificates for authentication. The certificates must specify both DNS names and IP addresses in the SAN field to prevent verification failures.

Configuration involves modifying `/var/packages/ContainerManager/etc/dockerd.json` to include both the unix socket (maintaining GUI functionality) and TCP endpoint with TLS settings. The configuration must specify certificate paths, typically stored in `/volume1/docker/certs/` with restrictive permissions (400 for private keys). After configuration changes, the Container Manager service requires a restart using `sudo systemctl restart pkg-ContainerManager-dockerd`.

Client systems need copies of the CA certificate, client certificate, and client key to authenticate. Connection methods include Docker contexts (`docker context create synology-tls --docker "host=tcp://nas-ip:2376,ca=ca.pem,cert=cert.pem,key=key.pem"`), environment variables (DOCKER_HOST, DOCKER_TLS_VERIFY, DOCKER_CERT_PATH), or command-line parameters. Certificate expiration represents a common issue, requiring regular rotation and monitoring to maintain connectivity.

## Native APIs offer limited container management capabilities

Synology's Web API provides programmatic access through endpoints like `SYNO.Docker.Container` (version 1), despite the rebrand to Container Manager. Authentication uses the DSM Login Web API to obtain session IDs, followed by container operations through `/webapi/entry.cgi`. The synowebapi CLI tool enables command-line API access, though it requires root privileges and generates deprecation warnings.

API functionality remains **limited compared to Docker's native capabilities**, supporting basic container lifecycle operations (start, stop, restart) but lacking advanced features like compose deployments or complex network configurations. Version compatibility issues exist between DSM versions, with documented version 7 endpoints sometimes less stable than version 6 equivalents. Public IP access often fails with permission errors, requiring local network connections or VPN tunnels for reliable operation.

Third-party tools like Portainer provide superior functionality, offering 2-3x faster operations and full Docker feature access while maintaining compatibility with Container Manager's backend. Installation through Task Scheduler or direct Docker deployment enables web-based management with better compose editing capabilities and comprehensive container orchestration features. Community projects on GitHub provide Synology-specific tools and workarounds for common limitations.

## Security requires multiple layers of protection

Production deployments demand **comprehensive security measures** beyond basic authentication. Network segmentation through Docker bridge networks isolates containers, while firewall rules must explicitly allow the Docker subnet (172.16.0.0/12) for proper container communication. The firewall rule order proves critical: local network allow rules must precede Docker subnet rules, which must appear before any deny rules.

User permission models range from full administrator access (docker group membership equals root privileges) to limited service accounts with specific UIDs mapped in compose files. The principle of least privilege suggests creating dedicated users for each container with minimal permissions, using the `--user` parameter or compose file specifications. Security options like `--security-opt no-new-privileges` and `--read-only` filesystems further restrict container capabilities.

VPN integration provides an additional security layer, with options including NAS-level VPN servers (OpenVPN or WireGuard), container-level VPN routing through GlueTUN, or mesh networking solutions like Tailscale. VPN usage may cause container connectivity issues requiring proper routing rules or split-tunnel configurations. Combining VPN access with SSH or TLS authentication creates defense-in-depth security suitable for production environments.

## Troubleshooting focuses on permissions and paths

Common issues center on **permission denied errors** and socket connection problems, typically resolved through proper group membership and socket ownership. Users must belong to the docker group (created with `sudo synogroup --add docker username`) and the socket requires docker group ownership with 660 permissions. Container Manager service issues after DSM updates often require package reinstallation or configuration file restoration from backups.

Certificate verification failures stem from incorrect system time, expired certificates, or path misconfigurations in dockerd.json. DSM's NTP synchronization through Control Panel prevents time-related issues, while certificate regeneration with proper validity periods and SAN entries resolves most TLS problems. Key length compatibility requires 2048-bit keys for some DSM configurations rather than the more secure 4096-bit default.

Performance issues manifest as slow container operations or network connectivity problems. Bridge network issues resolve through Docker service restarts, while persistent performance degradation may indicate insufficient resources or disk space. Monitoring through `docker stats`, system resource checks, and log analysis identifies bottlenecks. Container Manager's GUI performs significantly slower than CLI operations, making remote management essential for efficiency.

## Conclusion

Remote container management on Synology DSM 7.x requires understanding Container Manager's architectural changes and implementing appropriate workarounds for PATH and permission issues. SSH with Docker contexts provides the optimal balance of security and usability for most scenarios, while TLS-secured TCP connections suit automation needs. The critical PATH environment fix and proper certificate management enable reliable remote operations, with comprehensive security measures including VPN integration and network segmentation ensuring production-ready deployments. Regular monitoring of DSM updates remains essential as Synology continues evolving Container Manager's capabilities and potentially introducing breaking changes to established workarounds.