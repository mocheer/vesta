# Vesta

## docker 字体安装

### 命令行拷贝
本机为windows，字体目录为C:\Windows\Fonts，有很多系统是直接用window的本地字体，如果在linux部署，需要在linux环境中按照字体

- docker cp MSYH.ttc {Container_name_or_ID}:/usr/share/fonts/

### Dockerfile安装
本机为windows

```Dockerfile
# 解决截图中文乱码 其他字体请自行修改字体文件名
COPY MSYH.ttc /usr/share/fonts/ 
```


