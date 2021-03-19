<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [xterm web attach container](#xterm-web-attach-container)
  - [dev](#dev)
  - [problem](#problem)
  - [acknowledgement](#acknowledgement)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## xterm web attach container

xterm.js 连接 k8s 中某个 pod 的某个 container

### dev

- remotecommand
- xterm.js

### problem

- 中文问题：容器镜像本身需要支持中文， env LANG=C.UTF-8

```
kubectl exec -it nginx-deployment-66b6c48dd5-9bw5x -- env LANG=C.UTF-8 /bin/bash 
```

- Invalid UTF-8 in text frame, 使用 base64 先编码，传到前端再解码
- atob 原生不支持中文解码, 使用 webtoolkit Base64

### acknowledgement

- [容器中文问题](https://cloud.tencent.com/developer/article/1500399)
- [Invalid UTF-8](https://www.lflxp.cn/post/golang/websocket%E5%AE%9E%E6%88%98%E5%9B%9B/#websocket%E5%8D%8F%E8%AE%AE%E7%9A%84%E8%AF%B8%E5%A4%9A%E9%97%AE%E9%A2%98%E5%88%86%E6%9E%90%E5%92%8C%E8%A7%A3%E5%86%B3)
