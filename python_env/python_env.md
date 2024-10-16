## 在mac下使用不同版本的env



### 生成env环境

```bash
# 生成环境
python3 -m venv <path>/python3.12_env

# 激活环境
source <path>/python3.12_env/bin/activate
# 激活后在行首会有 (python3.12_env) 信息

# 安装包文件
pip3 install aa  && pip3 install -r reqxxx.txt
# 安装过后的包会在该目录下 <path>/python3.12_env/lib/python3.12/site-packages

# 取消激活
deactivate
```



