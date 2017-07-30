# dicam 0.5.0

dicam aims too allow easy management of a network of webcams for surveillance
or streaming.

**It is currently in beta, so fell free to report any bug or missing infos.**

## Features

- Motion detection from local or remote cameras
- Push notifications or email on motion
- Remote access and control via CLI
- HTML interface for webcams streaming and control

## Next steps

- [ ] User friendly notifications (_with images/videos if possible_)
- [ ] External feeds access (_dicam will act as a reverse proxy_)
- [ ] Web server index and stats/logs for each cameras
- [ ] Secure motion feed streaming
- [ ] Config file improvements
- [ ] Unit tests
- [ ] Live countdown cancelation (_via web interface or CLI_)
- [ ] Motion HTTP API support
- …

## Setup

### 1. Download dicam

dicam is available as a single binary for some platforms.

In order to get the last version, grab it from the latest
[release](https://github.com/davidderus/dicam/releases).

### 2. Run it!

#### A. Setup a config file

A dicam config file, stored in `~/.config/dicam/config.toml` may look like this:

```toml
countdown = 30 # In seconds, default to 10. Once the countdown is completed, an alert is sent.

# cameras are local or remote camera devices, which will be used for monitoring
# and streaming.
[cameras]
  [cameras.livingroom]
  auto_start = false # whether to start the camera on dicam's launch or not
  device = "/dev/video0" # a local video device path
  input = -1 # Custom video input handling for the given device (optional - default is 0)
  role = "watch" # Notifications and streaming.

  [cameras.dorms]
  auto_start = false # whether to start the camera on dicam's launch or not
  remote_device = "http://192.168.0.30:8080/feed.jpg" # a remote video device
  # path. It must be a URL that returns single jpeg pictures or a raw mjpeg stream
  role = "stream" # No notifications. Only streaming.

# the webserver allows cameras control (start/stop) and live preview.
[webserver]
  # optional HTTP digest auth
  [[webserver.user]]
    name = "john"
    password = "my_hashed_password"

# notifiers are used to alert about an event (a motion for example)
[notifiers]
  # Emailer is an email notifier, sending an email via a given SMTP server
  [notifiers.emailer]
    service = "email"
    recipients = [ "test@test.com" ]
    [notifiers.emailer.options]
      host = "my_host"
      port = "587" # must be a string
      from = "myemail@email.fr"
      password = "azerty"
  # pushbullet is a push notifier, pushing a message and/or picture to all
  # the recipients listed
  [notifiers.pushbullet]
    service = "push"
    recipients = [ "MYDEVICEID" ]
    [notifiers.pushbullet.options]
      api_key = "my_api_key"
```

#### B. Launch the controller (and command center)

In a shell (*via tmux or screen*), run the following command to launch the command center:

`dicam controller`

Once this is done, you can communicate with dicam with the cli and interact with
cameras and services.

The host and port can be changed in the config file as follow:

```toml
host = 0.0.0.0
port = 4541

# rest of the config
```

#### C. Interact with cameras

```shell
# Launch a camera
dicam cam start livingroom

# Stop a camera
dicam cam stop livingroom

# List all cams
dicam cam list
```

#### D. Live streams

Dicam embeds a webserver which allows you to control your cameras and watch their
streams in real time.

To start the webserver, run `dicam webserver`.

Then, you can access [0.0.0.0:8000]() and browse your cameras through the
web interface.

## How-to…

### Install motion

dicam leverages `motion` for webcams management and motion detection, so we need
the motion binary to be installed first.

*It was tested with motion `4.0.1` on macOS and with an earlier version
on Ubuntu.*

#### On Debian/Ubuntu

```
sudo apt install motion
```

#### On OS X

**[Only works with network webcams](https://github.com/jogu/motion/blob/master/README.MacOSX)**

```
# From https://github.com/Motion-Project/motion/wiki/cmake#macos-builds

brew install ffmpeg cmake pkg-config libjpeg
cmake .
make

cp motion /usr/local/bin/motion
```

### Generate a Web Interface password

The web interface supports HTTP Digest Authentication. The dicam realm is
`dicam.local` so use the following command to get your hashed password:

```shell
echo -n "$username:dicam.local:$password" | md5sum > .my-password
```

## Inner working

TODO
