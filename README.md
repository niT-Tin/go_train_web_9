# Vue+Go实现分布式火车购票系统
相关技术
1. 后台采用分布式微服务，服务内部使用grpc调用。
2. 使用Nacos作为后台配置中心。
3. 使用Jaeger实现关键接口的链路追踪。
3. 使用Sentinel实现对关键接口的流量限制。
4. 使用RocketMQ实现基于消息的数据最终一致性方案。
5. 使用前后端双验证码对流量进行削峰。

TODO:

- [ ] 相关流程图绘制以及文档编写(按道理这个因该最先弄😅)
- [ ] 实现自动化部署，编写DockerFile等
- [ ] 配置网关
