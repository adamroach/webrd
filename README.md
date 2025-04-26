# webrd
This is a very basic Remote Desktop server and client that allows remote desktop access via a web browser, using WebRTC. It currently works with MacOS desktops, but it should be straightforward to add Linux and Windows support.

It's still very rough around the edges, needs more testing, and needs documentation. 

**Don't turn off TLS unless you're really sure that you want everyone on the local network to spy on your keystrokes.**

# TODO
In no particular order:

- fix mouse pointer alignment
- documentation
- convert logging to zap
- audio capture
- bake client in with go:embed (make configurable?)
- macOS touchpad handling
- unit tests
- multi-screen support
- handling of oversized screens
- Windows support
- Linux support (Wayland)
- Linux support (X)