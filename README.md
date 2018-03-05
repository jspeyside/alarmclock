# alarmclock
Wake on LAN REST API

# Background
The original inspiration for this project was the desire to wake and shutdown a PC using SmartHome voice control. By running this service on your home network and adding a webhook to IFTTT, a user can use voice control to wake up and shutdown a machine.

# Usage
## Configuration
In order to successfully use alarmclock, a `config.yml` file must be created and completed. The configuration file has the following format:
```yaml
broadcast: "<broadcast_address>"
hosts:
  host1:
    mac: "12:34:56:78:90:AB"
    username: "<username>"
    password: "<password>"
  host2:
    mac: "12:34:56:78:90:AB"
    username: "<username>"
    password: "<password>"
```
The broadcast address for your home network can be obtained by running `ipconfig` Windows or `ifconfig` Linux/Mac from a terminal. Typically, this will be something similar to: `192.168.1.255`. The mac address for the machines to control wake/sleep can be found using the same commands.

Currently only a single broadcast is supported for all hosts. Also, only Windows support is available for shutdown. Username and password are not required in `config.yml` if shutdown functionality is not needed.

## Running Alarmclock
To run alarmclock, pass the location of the config file to the binary via:
```bash
./alarmclock -f <path_to_config>/config.yml
```

## Configuring Shutdown
### Windows
Windows needs special configuration to allow remote management and shutdown. To enable shutdown follow the instructions below.

On the remote host, a PowerShell prompt, using the **Run as Administrator** option and paste in the following lines:

```powershell
winrm quickconfig
y
winrm set winrm/config/service/Auth '@{Basic="true"}'
winrm set winrm/config/service '@{AllowUnencrypted="true"}'
winrm set winrm/config/winrs '@{MaxMemoryPerShellMB="1024"}'
```
Windows Firewall must be running to enable remote management. In addition the NIC used for shutdown, must be set as a private network (network discovery enabled).
