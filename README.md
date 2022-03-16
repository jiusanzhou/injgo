<div align="center">

# `injgo`

**Injgo** is a tool for dynamic library injecting which written in Golang.


[![](https://img.shields.io/travis/jiusanzhou/injgo.svg?label=build)](https://travis-ci.org/jiusanzhou/injgo) [![](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jiusanzhou/injgo) [![](https://goreportcard.com/badge/github.com/jiusanzhou/injgo)](https://goreportcard.com/report/jiusanzhou/injgo) [![@Zoe on Twitter](https://img.shields.io/badge/twitter-@jiusanzhou-55acee.svg)](https://twitter.com/jiusanzhou "@Zoe on Twitter") [![InjGo on Sourcegraph](https://sourcegraph.com/github.com/jiusanzhou/injgo/-/badge.svg)](https://sourcegraph.com/github.com/jiusanzhou/injgo?badge "InjGo on Sourcegraph")

</div>

|If you are a rustacean, try the Rust version: [injrs 🦀](https://github.com/jiusanzhou/injrs)|
|:---|

### Features

- **Pure `Go`**
- **Zero dependency**
- **Simple usage**

### Usage

You can use `injgo` as a cli tool.

**1. Install**

```bash
go get go.zoe.im/injgo/cmd/...
```

**2. Inject**

```bash
injgo PROCESS_NAME/PROCESS_ID DLL...
```

Also, you can use `injgo` as library.

### API

- `Inject(pid int, dllname string, replace bool) error`
- `InjectByProcessName(name string, dll string, replace bool) error`

### TODO

- [ ] Use injector to handle result
- [ ] Unload injected DLLs
