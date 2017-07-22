# dicam 0.5.0

dicam aims too allow easy management of a network of webcams for surveillance or streaming.

## Features

- Motion detection
- Push notifications or email on motion
- Remote access and control via CLI
- HTML interface for webcams streaming and control

## TODO

### For v1

- [ ] Handle IP cameras
- [ ] Bundle static assets in the binary
- [ ] Update documentation with examples, explanations and screenshots

### For later

- [ ] External feed access
- [ ] Config file improvements
- [ ] Web server index and stats/logs for each cameras
- [ ] Secure motion feed streaming
- [ ] Add some tests
- …

## Setup

### 1. Install motion

dicam leverages `motion` for webcams management and motion detection, so we need the motion binary to be installed first.

#### On Debian

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

### 2. Download dicam

dicam is available as a single binary for some platforms.

In order to get the last version, grab it from the latest [release](https://github.com/davidderus/dicam/releases).

### 3. Run it!

#### A. Setup a config file

A dicam config file, stored in `~/.config/dicam/config.toml` may look like this:

```toml
countdown = 30 # In seconds, default to 10

[cameras]
  [cameras.livingroom]
  auto_start = false
  device = "/dev/video0"
  role = "watch" # Notifications and streaming.

  [cameras.dorms]
  auto_start = false
  device = "/dev/video1"
  input = -1 # Custom input handling
  role = "stream" # No notifications. Only streaming.

[webserver]
  # optionnal HTTP digest auth
  [[webserver.user]]
    name = "john"
    password = "my_hashed_password"

[notifiers]
  [notifiers.emailer]
    service = "email"
    recipients = [ "test@test.com" ]
    [notifiers.emailer.options]
      host = "my_host"
      port = "587" # must be a string
      from = "myemail@email.fr"
      password = "azerty"
```

#### B. Launch the controller (and command center)

In a shell (*via tmux or screen*), run the following command to launch the command center:

`dicam controller`

Once this is done, you can communicate with dicam with the cli and interact with cameras and services.

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

##### Web Interface authentication

The web interface supports HTTP Digest Authentication. The dicam realm is `dicam.local`
so use the following command to get your hashed password:

```shell
echo -n "$username:dicam.local:$password" | md5sum > .my-password
```

#### E. Add other notifiers

##### Pushbullet

```toml
[notifiers]
  [notifiers.pushbullet]
    service = "push"
    recipients = [ "MYDEVICEID" ]
    [notifiers.pushbullet.options]
      api_key = "my_api_key"
```

## Inner working

TODO
