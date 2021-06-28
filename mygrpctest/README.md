# 对比服务器返回结果与期望值是否一致

## 目录结构
```
    cases：获取输出文件的路径
    check:完成值的对比
    client：接收服务器返回值，序列化到outputs文件夹中
    outputs:存返回的结果
    proro:存.proto文件以及protoc编译后的两个文件
    reporter:将返回结果写回到文件中
    testcase:存储测试文件
    util：获取测试文件路径
```
## 运行
```
    outpus中的文件应该是最新的结果
    go build 
    mygrpctest.exe initpron "testcase/case1" "outputs/case1"
```