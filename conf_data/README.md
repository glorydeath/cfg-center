# 全局配置文件托管repo，实时更新。
* 可以通过http://cfg-center-ip:2120/conf 访问获取json序列化的配置数据(Chrome推荐安装JSONView插件)
* 配置文件为YAML格式。
* 更详细的使用方式文档参见：http://github.com/4paradigm/cfg-center
* 通过webhook实现触发reload，不需要大家关注。
 
# 2016-09-05更新
* **key名字统一为xx-yyy-zz的形式**
* HDFS、Hadoop、YARN相关配置统一存放在common节
