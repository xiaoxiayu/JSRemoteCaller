"# JSRemoteCaller" 
##  部署
1. 服务编译及运行：go build -o service && ./service -port=9090
2. 需要调用的模块中加入调用代码：<script src="fxqa-test.js"></script>
3. 调用初始化: <script>fxqa_init()</script>

## 调用
### 直接调用：
#### 1. 执行代码：curl -d "YWxlcnQlMjglMjdvayUyNyUyOSUzQmZ4cWFfbG9nJTI4MCUyQyUyMCUyN29rJTI3JTI5JTNC" http://127.0.0.1:9090/code
#### 2. 等待执行完成：curl http://127.0.0.1:9090/wait
#### 3. 获取执行结果：curl http://127.0.0.1:9090/log

### 编写JS函数，PY调用：
#### 1. 在 js 中写好代码，如[fxqa_api_testcode.js](./fxqa_api_testcode.js) 中的 demo0
#### 2. 调用 fxqa_run('demo0')
