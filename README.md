# Drone 飞书通知插件

## 使用配置

```yml
  - name: 发送飞书通知
    image: ydq1234/drone-lark
    pull: if-not-exists
    when:
      status:
        - success
        - failure
    settings:
      webhook: https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxx
      secret: xxxxxxxxxxxx
      debug: true
```

## 本地构建

如果是本地私有仓库构建，修改 `build.sh` 里面的内容，执行 `./build.sh` 即可

```bash
vi ./build.sh
# replace ydq1234 => your account or docker hub domain

#给 build.sh 添加运行权限
chmod +x ./build.sh

#可以指定tag，不指定默认为 latest
#export tag=xxx
./build.sh
```