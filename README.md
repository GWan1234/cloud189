# cloud189

封装天翼云盘网页接口实现命令行访问, 目前为开发阶段, 不足之处还请包涵, 本人将持续优化使用体验

## 命令列表

命令中云端路径均以`/`开头, `...`表示支持多参数, 全局参数`--config`指定配置文件路径，默认路径为`${HOME}/.config/cloud189/config.json`，例：`cloud189 --config /tmp/config.json ls {云盘路径}`

- 显示帮助: `cloud189 -h`
- 显示版本: `cloud189 version`
- 用户登录: `cloud189 login`
  - `cloud189 login` 扫码登录, 在浏览器中打开控制台中输出的链接, 使用移动端扫码登录
  - `cloud189 login -i {用户名} {密码}` 用户名密码登录
- 每日签到: `cloud189 sign` 未实现客户端签到获取空间, 仅实现抽奖获取空间
- 查看空间: `cloud189 df` 查看云盘空间的使用信息
- 文件夹创建: `cloud189 mkdir {云盘路径}` 支持多层级目录创建
- 文件上传: `cloud189 up {本地路径|http|fast...} {云盘路径}`，支持三种模式文件上传，已知网页版上传接口存在bug不支持断点续传, 例
  - 本地上传`cloud189 up {本地路径...} {云盘路径}`，例 `cloud189 up /tmp/cloud189 /我的应用` 本地文件支持秒传
  - http上传 `cloud189 up {http://文件...} {云盘路径}`，例 `cloud189 up https://github.com/gowsp/cloud189/releases/download/v0.4.2/cloud189_0.4.2_linux_amd64.tar.gz /我的应用`，该模式不支持10M以上的文件秒传
  - 手动秒传 `cloud189 up {fast://文件MD5:文件大小/文件名...} {云盘路径}`，例 `cloud189 up fast://3BACAB45A36BE381390035D228BB23E0:7598080/cloud189 /我的应用`，可以实现无文件上传，例如：系统镜像
- 文件下载: `cloud189 dl {云端路径...} {本地路径}` 支持文件夹, 支持断点续传
- 文件列表: `cloud189 ls {云盘路径}` 大小为`-`表示文件夹
- 文件删除: `cloud189 rm {云盘路径...}`
- 文件复制: `cloud189 mv {云盘路径...} {目标路径}`
- 文件移动: `cloud189 cp {云盘路径...} {目标路径}`
- WebDAV: `cloud189 webdav :{端口}` 启动 webdav服务, 上传不支持10M以上的文件秒传
