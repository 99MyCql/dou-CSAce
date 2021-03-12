# douCSAce

本科毕业设计

题目：科研快速入门辅助系统

主要任务：为辅助科研人员快速入门，设计一个论文数据系统，能够基于论文不同研究方向呈现学术研究热点趋势以及相关方向的重要论文、学者信息，并能可视化呈现。

特色功能：
- 主要针对计算机领域的学术数据进行可视化分析与展示。
- 结合国内需求，根据《CCF推荐会议和期刊》进行研究方向的划分，可视化显示每个方向的变化趋势和对比。
- 在每个方向中，可视化显示其中每个会议/期刊论文数、被引用数的变化趋势、TOP作者论文信息，以及不同会议/期刊之间的对比。
- 在作者信息中，展示作者的合作关系网。

## 国内外情况

- 传统论文数据系统

![1](img/2.png)

- 新型学术图谱系统

![1](img/3.png)

## 数据模型设计

系统考虑使用知识图谱技术，知识图谱有三种存储方式：三元组、关系数据库、图数据库。关于数据库选型和 ArangoDB 数据库介绍，请参考我的博客：[ArangoDB入门](https://99mycql.github.io/application/ArangoDB%E5%85%A5%E9%97%A8.html) 。

模型按照 ArangoDB 数据的风格进行设计，分为实体模型和关系模型。

模型关系图如下：

![1](img/1.png)

各模型详细设计如下：

- 论文实体模型（Document），表名：`papers`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，等同于 dblp 中文章的 key ，`_key = <type>-<journal _key/conf _key>-<dblpkey>`，比如：`1-ppopp-GavrielatosKNGJ20`|
|title|string||
|abstract|string||
|type|uint|0:Other, 1:Conference, 2:Journal|
|pages|string||
|year|string||
|bookTitle|string|type=1|
|volume|string|type=2|
|number|string|type=2|
|doiUrl|string||
|dblpUrl|string||
|referenceCount|long||
|citationCount|long||

- 会议实例（哪一年哪一场的会议）实体模型（Document），表名：`conferenceInstances`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，等同于 dblp 中会议实例的 key ，比如：`ppopp2020`|
|title|string||
|publisher|string|出版方|
|booktitle|string||
|year|string||
|isbn|string|dblp 中会议实例的属性，用处未知，暂且保存下来|
|location|string|举办地点|
|doiUrl|string||
|dblpUrl|string||
|paperCount|uint||
|citationCount|uint||

- 会议实体模型（Document），表名：`conferenceSeries`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，等同于 dblp 中会议的 key ，`_key = shortName`，比如：`ppopp`|
|shortName|string||
|name|string||
|publisher|string||
|dblpUrl|string||
|paperCount|uint||
|citationCount|uint||

- 期刊实体模型（Document），表名：`journals`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，等同于 dblp 中期刊的 key ，`_key = shortName`，比如：`tocs`|
|shortName|string||
|name|string||
|publisher|string||
|dblpUrl|string||
|paperCount|uint||
|citationCount|uint||

- 研究方向实体模型（Document），表名：`fields`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，`_key = <type>-<name.replaceAll(" ", "_")>`|
|name|string||
|zhName|string||
|type|uint|0:Other, 1:CCF|
|paperCount|uint||
|citationCount|uint||

- 作者实体模型（Document），表名：`authors`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，等同于 dblp 中作者的 pid ，`_key = pid.replaceAll("/", "_")`，比如：`g_RajivGupta`|
|name|string||
|zhName|string||
|urls|string|作者主页，如果有多个则用空格分隔|
|paperCount|uint||
|citationCount|uint||

- 机构实体模型（Document），表名：`affiliations`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|name|string||
|paperCount|uint||
|authorCount|uint||
|citationCount|uint||

- 论文->论文 引用关系模型（Edge），表名：`cit_by`，From 表：`papers`，To 表：`papers`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中某个文档的 _key |
|_to|string| To 表中某个文档的 _key |

- 论文->会议实例 发表关系模型（Edge），表名：`publish_on_confIns`，From 表：`papers`，To 表：`conferenceInstances`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中某个文档的 _key |
|_to|string| To 表中某个文档的 _key |

- 论文->期刊 发表关系模型（Edge），表名：`publish_on_jou`，From 表：`papers`，To 表：`journals`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中某个文档的 _key |
|_to|string| To 表中某个文档的 _key |

- 会议实例->会议 从属关系模型（Edge），表名：`confIns_belong_to_confSer`，From 表：`conferenceInstances`，To 表：`conferenceSeries`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中某个文档的 _key |
|_to|string| To 表中某个文档的 _key |

- 会议->研究方向 从属关系模型（Edge），表名：`confSer_belong_to_field`，From 表：`conferenceSeries`，To 表：`fields`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中某个文档的 _key |
|_to|string| To 表中某个文档的 _key |
|note|string|备注，若指向的研究方向类型是 CCF ，则备注分类：A、B、C|

- 期刊->研究方向 从属关系模型（Edge），表名：`jou_belong_to_field`，From 表：`journals`，To 表：`fields`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中某个文档的 _key |
|_to|string| To 表中某个文档的 _key |
|note|string|备注，若指向的研究方向类型是 CCF ，则备注分类：A 类、B 类、C 类|

- 作者->作者 合作关系模型（Edge），表名：`co_with`，From 表：`authors`，To 表：`authors`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中对应文档的 _key |
|_to|string| To 表中对应文档的 _key |

- 作者->机构 从属关系模型（Edge），表名：`author_belong_to_aff`，From 表：`authors`，To 表：`affiliations`

|Field Name|Field Type|Description|
|---|---|---|
|_key|string|唯一标识，自动生成|
|_from|string| From 表中对应文档的 _key |
|_to|string| To 表中对应文档的 _key |
|startYear|string||
|endYear|string||
|note|string|备注|
