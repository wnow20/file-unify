# File Unify
 - unify - Unify files' encoding from source charset to target. 批量转换文件编码
 - detect - Detect file's charset. 检测文件编码

## How to use
### unify 文件编码转换
```shell
# unify files in current path
unify -s GB-18030 -t GBK -x test4gbk.txt
# unify files in a specified path
unify -s GB-18030 -t UTF8 -x .java -p /Users/ge/workspace/bc/businesscenter
```
### unify 文件编码检测
```shell
# 支持相对路径、绝对路径
./unify detect -f cmd/test4gbk.txt
```

## How to install
```shell
GOPROXY=https://goproxy.cn/,direct go install github.com/wnow20/unify@latest
```
