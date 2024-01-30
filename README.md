# Srun Login

西电校园网登录助手，通过 Web Portal 认证方式。可以运行在后台保持网络连接。支持自动重连、开机自启。

## 使用说明：

### 下载使用

- 在 Release 页面下载 zip 文件并解压
- 修改 `config.yaml`，填入学号和密码
- 运行 `srun-login.exe` 即可。（程序会在系统托盘后台运行）
- （可选）右键托盘图标，选择 “AutoStart”，即可开机自启

备注：

- 对于其他使用深澜 srun 校园网系统的高校，可以在 `config.yaml` 修改 `host` 来使用。
    ```
    host: http://domain
    ```
- 如果遇到问题，可以查看程序同目录下的日志文件 `log.txt` 来排查。

### 编译使用

编译程序：
```bash
go build -o srun-login.exe -ldflags "-s -w -H=windowsgui" ./app
```

在程序同目录创建 `config.yaml` 文件，写入登录信息：

```yaml
username: 23000000000
password: xxxxxx
```

运行程序即可。（程序会在系统托盘后台运行）

## 鸣谢

- [某校园网认证api分析 [ Trailblazer ]](https://www.ciduid.top/2022/0706/school-network-auth/)
- https://github.com/Debuffxb/srun-go