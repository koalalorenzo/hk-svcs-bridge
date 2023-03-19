# Go HomeKit Services Bridge
o-hk-systemd is a system service written in Go that acts as a HomeKit bridge to
connect SystemD services and HomeKit automation. This project is available for
macOS and Linux, and it allows users to add new fake switches to their Apple
HomeKit home setup. When these switches are turned on, they will start some
SystemD services, and when turned off, they will stop them.

This project is designed to be easy to use and customizable. It supports custom
commands, so it does not have to be integrated with SystemD. Additionally, it
can be used with any SystemD service, making it a flexible and versatile tool
for managing services.

**Important**: This is a WIP and side project. It is not designed for system 
production usage and this project is just for fun!

## How to install from source
If you are using GNU/Linux with SystemD you can automatically install the binary 
from the source code. Make sure to have installed:

* SystemD (available by default on many distro like Ubuntu, Debian and Arch)
* GNU Make
* Go
* Git

You can then run these commands to install the binary and the service:

```bash
git clone https://gitlab.com/koalalorenzo/go-hk-systemd-bridge.git
cd go-hk-systemd-bridge
sudo make install
```

After that, you **must change the configuration** in 
`/etc/go-hk-systemd-bridge.yaml`. Then you can run/restart the bridge by 
running:

```bash
sudo systemd restart go-hk-systemd-bridge
```

## Configuration
You can customize many things, from the name of the bridge in Home Kit to 
the services and commands that will be executed. You can read more
in the example configuration file [here](config.yaml).

Please note that many of the options can be omitted, the default values are 
specified in the [code for the configuration](config.go).

By default the binary will look for `config.yaml` path, but you can customize
the path by using the environment variable `CONFIG`:

```bash
export CONFIG=/etc/go-hk-systemd-bridge.yaml
/usr/bin/go-hk-systemd-bridge
```

