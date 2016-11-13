# dicam 0.3.0

dicam aims too allow easy management of a network of webcams for surveillance or streaming.

## Features

- Motion detection
- Push, email and sound alerts on motion
- Remote access and control via CLI
- HTML Webserver for webcams streaming
- Authentificated activation and deactivation

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
cp motion-dist.conf /usr/local/etc/motion/motion.conf
```

### 2. Download dicam

dicam is available as a single binary for some platforms.

In order to get the last version, grab it from the latest [release](https://github.com/davidderus/dicam/releases).

### 3. Run it!

TODO

## Inner working

TODO
