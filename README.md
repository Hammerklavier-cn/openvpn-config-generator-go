# openvpn-config-generator-go

[English](README_en#openvpn-config-generator-go)

**警告：该项目还在开发中，请不要使用。**  
**Warning: This project is still under development. Please do not use it.**

自动生成 OpenVPN 的证书文件。Automatically generate configurations for OpenVPN.

## 依赖

1. 通过包管理器 `easy-rsa` 和 `openvpn`。
2. 程序运行需要 `sudo` 权限

## 程序运行

首先在你希望存放证书的位置运行 `opvn-gen init`，再在相同目录运行 `opvn-gen client` 获取证书

```text
Usage:
  opvn-gen [flags]
  opvn-gen [command]

Available Commands:
  client      Generate client certificate (.opvn).
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialise configuration and creates server configuration.

Flags:
  -h, --help      help for opvn-gen
  -V, --verbose   Show verbose process.
  -v, --version   version for opvn-gen

Use "opvn-gen [command] --help" for more information about a command.
```

`opvn-gen init` 用法：

```text
Usage:
  opvn-gen init [flags]

Flags:
      --algorithm string   Set algorithm for certificate. (default "RSA")
  -d, --days int           For how long the certificate remains valid. (default 180)
  -h, --help               help for init
      --keysize int        Set key size. (default 2048)
  -p, --path string        Set the directory where configuration files are stored. (default ".")

Global Flags:
  -V, --verbose   Show verbose process.
```

`opvn-gen client` 用法：

```text
Usage:
  opvn-gen client [flags]

Flags:
  -h, --help          help for client
  -n, --name string   Set the client name. (default "client_cert")
  -p, --path string   Set the directory where client configuration files are stored. (default ".")

Global Flags:
  -V, --verbose   Show verbose process.
```

## 源码编译

只需要最新的 Go 编译器即可。撰写程序所用 Go 编译器版本为 1.23.6。
