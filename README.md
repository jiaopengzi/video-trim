# video-trim

video-trim 视频裁剪的意思。开发这个工具的初衷是，在手机上下载的视频(你懂得)，总是有一些开头一样的背景，导致手机打开相册的时候看到的所有都是一样的封面图，想找个自己要的视频，找半天都找不到，这不能忍。

最初我使用手机上相册自带的剪辑功能，盲猜是因为要编码所以慢的要死。

我是想要快速的将指定时间段的剪切掉即可，这太浪费时间了，也想过自己在手机上使用 ffmpeg 来，使用命令行来实现，但是要 root 权限，还是放弃了。

所以就将这个功能转移到了自己局域网的工作电脑了，性能跟得上，节约时间。所以就套壳一下 ffmpeg，手机使用浏览器访问即可实现使用。

## 思路

使用 golang 写一个 web 服务，在任意浏览器上实现选取本地的视频文件，然后实现批量的裁剪，裁剪好的视频，可以批量下载。

- 内网上传下载飞快
- 只是裁剪编码，使用 ffmpeg 的 copy 参数即可，飞快

## 用法

### 依赖安装

本工具依赖 ffmpeg, ffprobe，请事先安装好。 通过 `ffmpeg -version` `ffprobe -version` 查看是否安装成功。

windows 安装 ffmpeg 在 powershell 中运行 `scoop install ffmpeg`

其他系统参考 <https://www.gyan.dev/ffmpeg/builds/>

### 工具下载

在[releases](https://github.com/jiaopengzi/video-trim/releases)下载对应的版本。

将文件解压后直接双击 `video-trim-windows.exe` 即可。

然后根据提示访问对应的链接即可快乐裁剪视频了。
