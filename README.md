# vesta

一个网页采集的相关类库

### linux 安装 chrome

```sh
yum install mesa-libOSMesa-devel gnu-free-sans-fonts wqy-zenhei-fonts
yum install https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm 
yum install ./google-chrome-stable_current_x86_64.rpm

# 查看chrome版本号
google-chrome --version
```

### docker 
```dockerfile
# 指定基础镜像
FROM chromedp/headless-shell:latest
```