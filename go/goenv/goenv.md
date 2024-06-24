# Goenv使用



Goenv是一个**用于管理Go语言版本的工具**。 它允许您在同一台计算机上同时安装和切换不同的Go语言版本



## 安装（Mac）

```bash
# 安装
brew install goenv 

# 设置环境变量
echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.bash_profile
echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bash_profile
echo 'eval "$(goenv init -)"' >> ~/.bash_profile

```



## go版本管理

```bash
# 查看可用go版本
goenv list -l

# 安装指定版本
goenv install 1.20 
goenv install 1.21

# 查看已安装的版本
goenv versions

# 删除/卸载 指定版本
goenv uninstall 1.20 

```



## 设置goenv全局版本

```bash
# 指定全局go版本
goenv global 1.21 
```



## 设置指定项目goenv版本

```bash
# 指定项目使用go版本
cd /path/to/your/project
goenv local 1.20

# 确认版本
goenv version 
```



## 对应版本外目录

```bash 
# go mod 对应 pkg
～/.goenv/versions

# 目录结构如下
～/.goenv/
├── bin
├── versions
│   ├── 1.20.14
│   │   ├── bin
│   │   ├── pkg
│   │   ├── src
│   │   └── ...
│   ├── 1.21.11
│   │   ├── bin
│   │   ├── pkg
│   │   ├── src
│   │   └── ...
│   └── ...
├── shims
├── plugins
└── ...

```

