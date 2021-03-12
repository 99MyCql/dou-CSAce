package pkg

import (
	"github.com/sirupsen/logrus"
)

// 全局常量
const (
	ProjectName = "douCSAce" // 项目名

	// 各模型名
	PaperName               = "paper"                     // 论文 实体模型
	ConferenceInsName       = "conferenceInstance"        // 会议实例 实体模型
	ConferenceSerName       = "conferenceSeries"          // 会议 实体模型
	JournalName             = "journal"                   // 期刊 实体模型
	FieldName               = "field"                     // 研究方向 实体模型
	AuthorName              = "author"                    // 作者 实体模型
	AffiliationName         = "affiliation"               // 机构 实体模型
	CitByName               = "cit_by"                    // 论文->论文 引用关系模型
	PublishOnConfInsName    = "publish_on_confIns"        // 论文->会议实例 发表关系模型
	PublishOnJouName        = "publish_on_jou"            // 论文->期刊 发表关系模型
	ConfInsBelToConfSerName = "confIns_belong_to_confSer" // 会议实例->会议 从属关系模型
	ConfSerBelToFieldName   = "confSer_belong_to_field"   // 会议->研究方向 从属关系模型
	JouBelToFieldName       = "jou_belong_to_field"       // 期刊->研究方向 从属关系模型
	CoWithName              = "co_with"                   // 作者->作者 合作关系模型
	AuthorBelToAffName      = "author_belong_to_aff"      // 作者->机构 从属关系模型
)

// 全局变量
var (
	Conf *Config        // 配置
	Log  *logrus.Logger // 日志
	DB   *DBInfo        // 数据库
)
