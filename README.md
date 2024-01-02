# Drone 飞书通知插件

## 使用配置

```yml
  - name: 发送飞书通知
    image: ydq1234/drone-lark:0.0.1
    pull: if-not-exists
    when:
      status:
        - success
        - failure
    settings:
      webhook: https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxx
      secret: xxxxxxxxxxxx
```

## 本地构建

如果是本地私有仓库构建，修改 `build.sh` 里面的内容，执行 `./build.sh` 即可

```bash
chmod +x ./build.sh
export tag=xxx
./build.sh
```