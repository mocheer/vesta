# Vesta

## docker 字体安装


### Dockerfile安装
本机为windows

```Dockerfile
# 解决截图中文乱码 其他字体请自行修改字体文件名
COPY MSYH.ttc /usr/share/fonts/ 
```

### 命令行拷贝
本机为windows，字体目录为C:\Windows\Fonts

- docker cp MSYH.ttc {Container_name_or_ID}:/usr/share/fonts/
