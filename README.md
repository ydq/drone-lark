# Drone Lark Plugin 飞书通知插件

## Sample

```yml
  - name: Send Lark Message
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

## Build

```bash
vi ./build.sh
# replace ydq1234 => your account or docker hub domain

chmod +x ./build.sh

#export tag=xxx  （tag default is latest）
./build.sh
```