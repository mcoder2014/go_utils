package sendcloud

import (
	"encoding/json"
	"net/url"
)

type NormalSendEmailParam struct {
	ApiUser         string `json:"apiUser"`
	ApiKey          string `json:"apiKey"`
	From            string `json:"from"`                      // 	是	发件人地址. 举例: support@ifaxin.com, 为了更高的送达率，建议from域名后缀与发信域名一致。
	To              string `json:"to"`                        // 	收件人地址. 多个地址使用';'分隔, 如 ben@ifaxin.com;joe@ifaxin.com
	Subject         string `json:"subject"`                   // 	标题. 不能为空
	HTML            string `json:"html"`                      // 	邮件的内容. 邮件格式为 text/html
	ContentSummary  string `json:"contentSummary,omitempty"`  // 	邮件摘要. 该字段传入值后，若原邮件已有摘要，会覆盖原邮件的摘要；若原邮件中没有摘要将会插入摘要。了解邮件摘要的更多内容，请点击这里
	FromName        string `json:"fromName,omitempty"`        // 	发件人名称. 显示如: ifaxin客服支持<support@ifaxin.com>
	CC              string `json:"cc,omitempty"`              // 	抄送地址. 多个地址使用';'分隔
	BCC             string `json:"bcc,omitempty"`             // 	密送地址. 多个地址使用';'分隔
	ReplyTo         string `json:"replyTo,omitempty"`         // 	设置用户默认的回复邮件地址.多个地址使用';'分隔，地址个数不能超过3个. 如果 replyTo 没有或者为空, 则默认的回复邮件地址为 from
	LabelName       string `json:"labelName,omitempty"`       // 	本次发送所使用的标签名称. 若此标签名称在账户不存在，将自动创建.
	SendTags        string `json:"sendTags,omitempty"`        // 	本次发送所使用的tag. 多个tag使用';'分隔，每个tag长度不能超过20字符，最多5个tag.
	Headers         string `json:"headers,omitempty"`         // 	邮件头部信息. JSON 格式, 比如:{"header1": "value1", "header2": "value2"}
	Attachments     string `json:"attachments,omitempty"`     // 	邮件附件. 发送附件时, 必须使用 multipart/form-data 进行 post 提交 (表单提交)
	XSMTPAPI        string `json:"xsmtpapi,omitempty"`        // 	SMTP 扩展字段. 如 {"to": ["support@ifaxin.com"],"sub":{}},详见 X-SMTPAPI
	Plain           string `json:"plain,omitempty"`           // 	邮件的内容. 邮件格式为 text/plain
	SendRequestId   string `json:"sendRequestId,omitempty"`   // 	最大支持128字符，1小时内相同sendRequestId的多次请求只会处理第一次。
	RespEmailId     bool   `json:"respEmailId,omitempty"`     // 	默认值: true. 是否返回 emailId. 有多个收件人时, 会返回 emailId 的列表
	UseNotification bool   `json:"useNotification,omitempty"` // 	默认值: false. 是否使用回执
	UseAddressList  bool   `json:"useAddressList,omitempty"`  // 	默认值: false. 是否使用地址列表发送. 比如: to=group1@maillist.sendcloud.org;group2@maillist.sendcloud.org
}

func (p *NormalSendEmailParam) ToURLValue() (*url.Values, error) {
	return StructToURLValues(p)
}

type BaseResp struct {
	Result     bool            `json:"result"`
	StatusCode int64           `json:"statusCode"`
	Message    string          `json:"message"`
	Info       json.RawMessage `json:"info,omitempty"` //不同接口不一样
}

type NormalPlanEmailParam struct {
	From     string `json:"from"`
	FromName string `json:"fromName,omitempty"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Plain    string `json:"plain"`

	SendRequestId  string `json:"sendRequestId,omitempty"`
	ContentSummary string `json:"contentSummary,omitempty"`
	CC             string `json:"cc,omitempty"`
	BCC            string `json:"bcc,omitempty"`
	ReplyTo        string `json:"replyTo,omitempty"`
}

func (p *NormalPlanEmailParam) ToNormalSendEmailParam(apiuser, apitoken string) *NormalSendEmailParam {
	return &NormalSendEmailParam{
		ApiUser:        apiuser,
		ApiKey:         apitoken,
		From:           p.From,
		FromName:       p.FromName,
		To:             p.To,
		Subject:        p.Subject,
		Plain:          p.Plain,
		SendRequestId:  p.SendRequestId,
		ContentSummary: p.ContentSummary,
		CC:             p.CC,
		BCC:            p.BCC,
		ReplyTo:        p.ReplyTo,
	}
}

type ErrCode int64

const CodeSuccess ErrCode = 200    //	请求成功
const ErrCode40000 ErrCode = 40000 //	重复请求
const ErrCode40001 ErrCode = 40001 //	start不能为空
const ErrCode40002 ErrCode = 40002 //	start非法
const ErrCode40003 ErrCode = 40003 //	limit不能为空
const ErrCode40004 ErrCode = 40004 //	limit非法
const ErrCode40005 ErrCode = 40005 //	认证失败
const ErrCode40006 ErrCode = 40006 //	days格式非法,必须是大于0的正整数
const ErrCode40007 ErrCode = 40007 //	startDate格式错误,应该类似'2013-03-19'
const ErrCode40008 ErrCode = 40008 //	endDate格式错误,应该类似'2013-03-19'
const ErrCode40009 ErrCode = 40009 //	labelIdList不能为空
const ErrCode40010 ErrCode = 40010 //	apiUserList不能为空
const ErrCode40011 ErrCode = 40011 //	email不能为空
const ErrCode40012 ErrCode = 40012 //	email格式非法
const ErrCode40013 ErrCode = 40013 //	domainList不能为空
const ErrCode40014 ErrCode = 40014 //	标签ID不能为空
const ErrCode40015 ErrCode = 40015 //	标签ID格式错误
const ErrCode40016 ErrCode = 40016 //	apiUserList格式非法
const ErrCode40017 ErrCode = 40017 //	聚合参数格式错误
const ErrCode40100 ErrCode = 40100 //	标签创建成功
const ErrCode40101 ErrCode = 40101 //	标签创建失败
const ErrCode40102 ErrCode = 40102 //	标签ID不能为空
const ErrCode40103 ErrCode = 40103 //	标签ID非法
const ErrCode40104 ErrCode = 40104 //	标签名称不能为空
const ErrCode40105 ErrCode = 40105 //	标签名称的长度应该为1-255个字符
const ErrCode40106 ErrCode = 40106 //	标签ID对应的标签不存在
const ErrCode40107 ErrCode = 40107 //	标签删除成功
const ErrCode40108 ErrCode = 40108 //	标签删除失败
const ErrCode40109 ErrCode = 40109 //	标签更新成功
const ErrCode40110 ErrCode = 40110 //	标签更新失败
const ErrCode40111 ErrCode = 40111 //	query不能为空
const ErrCode40112 ErrCode = 40112 //	query的长度的长度应该为1-255个字符
const ErrCode40113 ErrCode = 40113 //	标签名称已经存在
const ErrCode40201 ErrCode = 40201 //	模版调用名称invokeName不能为空
const ErrCode40202 ErrCode = 40202 //	模版调用名称invokeName格式错误
const ErrCode40203 ErrCode = 40203 //	模版类型不能为空
const ErrCode40204 ErrCode = 40204 //	非法的模板类型, 只能是0或者1
const ErrCode40205 ErrCode = 40205 //	templateStat不能为空
const ErrCode40206 ErrCode = 40206 //	templateStat非法, 只能是-1, -2, 1, 0中的值
const ErrCode40207 ErrCode = 40207 //	name不能为空
const ErrCode40208 ErrCode = 40208 //	name格式非法
const ErrCode40209 ErrCode = 40209 //	subject不能为空
const ErrCode40210 ErrCode = 40210 //	subject格式非法
const ErrCode40211 ErrCode = 40211 //	html不能为空
const ErrCode40212 ErrCode = 40212 //	html格式非法
const ErrCode40213 ErrCode = 40213 //	text不能为空
const ErrCode40214 ErrCode = 40214 //	text格式非法
const ErrCode40215 ErrCode = 40215 //	模板创建失败
const ErrCode40216 ErrCode = 40216 //	模板调用名称对应的模板不存在
const ErrCode40217 ErrCode = 40217 //	模板删除失败
const ErrCode40218 ErrCode = 40218 //	模板更新失败
const ErrCode40219 ErrCode = 40219 //	用户最多只能有50个模板
const ErrCode40220 ErrCode = 40220 //	模版调用名称已经存在
const ErrCode40221 ErrCode = 40221 //	isSubmitAudit不能为空
const ErrCode40222 ErrCode = 40222 //	isSubmitAudit格式错误
const ErrCode40223 ErrCode = 40223 //	模板处于待审核状态, 不能修改
const ErrCode40224 ErrCode = 40224 //	cancel不能为空
const ErrCode40225 ErrCode = 40225 //	cancel格式错误
const ErrCode40226 ErrCode = 40226 //	模板处于待审核状态, 无需再次提交
const ErrCode40227 ErrCode = 40227 //	模板已经审核通过, 无需再次提交
const ErrCode40228 ErrCode = 40228 //	模板处于审核失败状态, 无需撤销审核
const ErrCode40229 ErrCode = 40229 //	模板还未提交审核, 无法撤销审核
const ErrCode40230 ErrCode = 40230 //	模板调用名称 与 (开始日期 、结束日期) 参数两者不能同时为空
const ErrCode40231 ErrCode = 40231 //	taskId必须是数值
const ErrCode40232 ErrCode = 40232 //	taskName不能为空
const ErrCode40233 ErrCode = 40233 //	taskName长度超过256个字符
const ErrCode40234 ErrCode = 40234 //	runTime不能为空
const ErrCode40235 ErrCode = 40235 //	runTime格式不正确，需为yyyy-MM-dd HH:mm:ss
const ErrCode40236 ErrCode = 40236 //	runTime必须晚于当前时间
const ErrCode40237 ErrCode = 40237 //	未指定上传文件名
const ErrCode40238 ErrCode = 40238 //	文件名不能超过255个字符
const ErrCode40239 ErrCode = 40239 //	文件名后缀不能为空
const ErrCode40240 ErrCode = 40240 //	不能上传%s格式文件（如js,com,exe,sh,cs）
const ErrCode40241 ErrCode = 40241 //	存储空间不够，请删除部分文件后再上传
const ErrCode40242 ErrCode = 40242 //	定时任务创建成功
const ErrCode40243 ErrCode = 40243 //	定时任务创建失败
const ErrCode40244 ErrCode = 40244 //	定时任务删除成功
const ErrCode40245 ErrCode = 40245 //	定时任务删除失败
const ErrCode40246 ErrCode = 40246 //	定时任务更新成功
const ErrCode40247 ErrCode = 40247 //	定时任务更新失败
const ErrCode40248 ErrCode = 40248 //	文件上传失败！
const ErrCode40249 ErrCode = 40249 //	有邮件附件需要上传,runTime需在当前时间上延后5分钟！
const ErrCode40250 ErrCode = 40250 //	重复创建定时任务
const ErrCode40251 ErrCode = 40251 //	频繁更新定时任务
const ErrCode40252 ErrCode = 40252 //	taskId对应的定时任务不存在
const ErrCode40253 ErrCode = 40253 //	仅还未开始执行的定时任务可以删除
const ErrCode40254 ErrCode = 40254 //	仅还未开始执行的定时任务可以更新
const ErrCode40255 ErrCode = 40255 //	地址列表中成员地址的个数不能为0
const ErrCode40256 ErrCode = 40256 //	邮件附件上传超时
const ErrCode40257 ErrCode = 40257 //	发信人地址from最大128个字符
const ErrCode40258 ErrCode = 40258 //	参数emptyExistedAttachments不能为空,且值只能为0或1
const ErrCode40401 ErrCode = 40401 //	取消订阅记录创建成功
const ErrCode40402 ErrCode = 40402 //	取消订阅记录创建失败
const ErrCode40403 ErrCode = 40403 //	取消订阅记录删除成功
const ErrCode40404 ErrCode = 40404 //	取消订阅记录删除失败
const ErrCode40501 ErrCode = 40501 //	name不能为空
const ErrCode40502 ErrCode = 40502 //	地址列表名称的长度应该为1-48个字符
const ErrCode40503 ErrCode = 40503 //	address不能为空
const ErrCode40504 ErrCode = 40504 //	地址列表别名的长度应该为1-48个字符
const ErrCode40505 ErrCode = 40505 //	地址列表别名已经存在
const ErrCode40506 ErrCode = 40506 //	desc不能为空
const ErrCode40507 ErrCode = 40507 //	地址列表描述的长度应该为1-250个字符
const ErrCode40508 ErrCode = 40508 //	地址列表创建失败
const ErrCode40509 ErrCode = 40509 //	newAddress不能为空
const ErrCode40510 ErrCode = 40510 //	新的地址列表别名的长度应该为1-48个字符
const ErrCode40511 ErrCode = 40511 //	address参数错误
const ErrCode40512 ErrCode = 40512 //	members不能为空
const ErrCode40513 ErrCode = 40513 //	成员地址的长度应该为1-48个字符
const ErrCode40514 ErrCode = 40514 //	成员地址的个数不能小于0
const ErrCode40515 ErrCode = 40515 //	成员地址的个数不能超过1000
const ErrCode40516 ErrCode = 40516 //	添加成员失败
const ErrCode40517 ErrCode = 40517 //	地址列表不属于此用户
const ErrCode40518 ErrCode = 40518 //	成员地址不符合邮件地址规范
const ErrCode40519 ErrCode = 40519 //	删除成员失败
const ErrCode40520 ErrCode = 40520 //	vars不能为空
const ErrCode40521 ErrCode = 40521 //	vars参数中变量个数和成员地址个数不相等
const ErrCode40522 ErrCode = 40522 //	vars参数不符合JSON字符串语法
const ErrCode40531 ErrCode = 40531 //	单成员变量最大长度为1024个字符
const ErrCode40601 ErrCode = 40601 //	退信记录删除成功
const ErrCode40602 ErrCode = 40602 //	退信记录删除失败
const ErrCode40603 ErrCode = 40603 //	邮箱地址已经存在
const ErrCode40604 ErrCode = 40604 //	过期时间格式为： 2018-03-19
const ErrCode40701 ErrCode = 40701 //	分组ID不能为空
const ErrCode40702 ErrCode = 40702 //	分组ID格式错误
const ErrCode40703 ErrCode = 40703 //	事件类型不能为空
const ErrCode40704 ErrCode = 40704 //	事件类型格式错误,没有可用的事件类型
const ErrCode40705 ErrCode = 40705 //	url不能为空
const ErrCode40706 ErrCode = 40706 //	url格式错误
const ErrCode40707 ErrCode = 40707 //	url测试失败
const ErrCode40708 ErrCode = 40708 //	url已存在
const ErrCode40709 ErrCode = 40709 //	groupId对应的webhook配置未找到
const ErrCode40710 ErrCode = 40710 //	webhook配置创建失败
const ErrCode40711 ErrCode = 40711 //	webhook配置删除失败
const ErrCode40712 ErrCode = 40712 //	webhook配置修改失败
const ErrCode40801 ErrCode = 40801 //	发信人地址from不能为空
const ErrCode40802 ErrCode = 40802 //	发信人地址from格式错误
const ErrCode40803 ErrCode = 40803 //	发信人名称fromName不能为空
const ErrCode40804 ErrCode = 40804 //	发信人名称fromName格式错误
const ErrCode40805 ErrCode = 40805 //	收件人地址不能为空
const ErrCode40806 ErrCode = 40806 //	收件人地址数组中, 存在非法地址
const ErrCode40807 ErrCode = 40807 //	收件人地址的数目不能超过100
const ErrCode40808 ErrCode = 40808 //	邮件主题subject不能为空
const ErrCode40809 ErrCode = 40809 //	邮件主题subject格式错误
const ErrCode40810 ErrCode = 40810 //	回复地址replyto不能为空
const ErrCode40811 ErrCode = 40811 //	回复地址replyto格式错误
const ErrCode40812 ErrCode = 40812 //	xsmtpapi不能为空
const ErrCode40813 ErrCode = 40813 //	xsmtpapi格式错误
const ErrCode40814 ErrCode = 40814 //	xsmtpapi解析值不能为空
const ErrCode40815 ErrCode = 40815 //	xsmtpapi必须含有to字段
const ErrCode40816 ErrCode = 40816 //	xsmtpapi中to字段的解析值不能为空
const ErrCode40817 ErrCode = 40817 //	xsmtpapi解析错误
const ErrCode40818 ErrCode = 40818 //	attachments不能为空
const ErrCode40819 ErrCode = 40819 //	附件大小不能超过10485760字节
const ErrCode40820 ErrCode = 40820 //	此用户没有使用地址列表的权限
const ErrCode40821 ErrCode = 40821 //	地址列表任务创建成功
const ErrCode40822 ErrCode = 40822 //	地址列表任务创建失败
const ErrCode40823 ErrCode = 40823 //	邮件模板不存在
const ErrCode40824 ErrCode = 40824 //	模板未通过审核
const ErrCode40825 ErrCode = 40825 //	邮件模板和API_USER类型不匹配
const ErrCode40826 ErrCode = 40826 //	参数subject和模板主题不能同时为空
const ErrCode40827 ErrCode = 40827 //	xsmtpapi中to数组长度不能超过100
const ErrCode40828 ErrCode = 40828 //	回执地址不能为空
const ErrCode40829 ErrCode = 40829 //	回执地址格式错误
const ErrCode40830 ErrCode = 40830 //	plain内容不能为空
const ErrCode40831 ErrCode = 40831 //	plain内容格式错误
const ErrCode40832 ErrCode = 40832 //	会议起始时间startTime不能为空
const ErrCode40833 ErrCode = 40833 //	会议起始时间startTime格式错误
const ErrCode40834 ErrCode = 40834 //	会议结束时间endTime不能为空
const ErrCode40835 ErrCode = 40835 //	会议结束时间endTime格式错误
const ErrCode40836 ErrCode = 40836 //	会议标题title不能为空
const ErrCode40837 ErrCode = 40837 //	会议标题title格式错误
const ErrCode40838 ErrCode = 40838 //	会议组织者名称不能为空
const ErrCode40839 ErrCode = 40839 //	会议组织者名称格式错误
const ErrCode40840 ErrCode = 40840 //	会议组织者邮件地址不能为空
const ErrCode40841 ErrCode = 40841 //	会议组织者邮件地址格式错误
const ErrCode40842 ErrCode = 40842 //	会议地点location不能为空
const ErrCode40843 ErrCode = 40843 //	会议地点location格式错误
const ErrCode40844 ErrCode = 40844 //	会议描述description不能为空
const ErrCode40845 ErrCode = 40845 //	会议描述description格式错误
const ErrCode40846 ErrCode = 40846 //	会议参与者名称不能为空
const ErrCode40847 ErrCode = 40847 //	会议参与者名称格式不对
const ErrCode40848 ErrCode = 40848 //	会议参与者邮件地址不能为空
const ErrCode40849 ErrCode = 40849 //	会议参与者邮件地址格式错误
const ErrCode40850 ErrCode = 40850 //	会议参与者名称个数和会议参与者邮件地址个数不相等
const ErrCode40851 ErrCode = 40851 //	会议邮件拼装失败
const ErrCode40852 ErrCode = 40852 //	cc地址不能为空
const ErrCode40853 ErrCode = 40853 //	cc地址格式错误
const ErrCode40854 ErrCode = 40854 //	CC地址的数目不能超过100
const ErrCode40855 ErrCode = 40855 //	bcc地址不能为空
const ErrCode40856 ErrCode = 40856 //	bcc地址格式错误
const ErrCode40857 ErrCode = 40857 //	BCC地址的数目不能超过100
const ErrCode40858 ErrCode = 40858 //	respEmailId不能为空
const ErrCode40859 ErrCode = 40859 //	respEmailId格式错误
const ErrCode40860 ErrCode = 40860 //	gzipCompress不能为空
const ErrCode40861 ErrCode = 40861 //	gzipCompress格式错误
const ErrCode40862 ErrCode = 40862 //	to中有格式错误的地址列表
const ErrCode40863 ErrCode = 40863 //	to中有不存在的地址列表
const ErrCode40864 ErrCode = 40864 //	地址列表的数目不能超过5
const ErrCode40865 ErrCode = 40865 //	html解压失败
const ErrCode40866 ErrCode = 40866 //	plain解压失败
const ErrCode40867 ErrCode = 40867 //	处理附件发生异常
const ErrCode40868 ErrCode = 40868 //	headers不能为空
const ErrCode40869 ErrCode = 40869 //	headers格式错误
const ErrCode40870 ErrCode = 40870 //	html和plain不能同时为空
const ErrCode40871 ErrCode = 40871 //	html格式错误
const ErrCode40872 ErrCode = 40872 //	邮件列表地址不能为空
const ErrCode40873 ErrCode = 40873 //	useAddressList不能为空
const ErrCode40874 ErrCode = 40874 //	useAddressList格式错误
const ErrCode40875 ErrCode = 40875 //	内嵌图片ID或内嵌图片附件长度不一致
const ErrCode40876 ErrCode = 40876 //	是否取消日程参数isCancel格式错误
const ErrCode40877 ErrCode = 40877 //	摘要不能为空
const ErrCode40878 ErrCode = 40878 //	摘要长度不能超过255个字节
const ErrCode40879 ErrCode = 40879 //	回复地址replyto个数不能超过3个
const ErrCode40880 ErrCode = 40880 //	xsmtpapi中to字段含有非法邮箱格式
const ErrCode40901 ErrCode = 40901 //	邮件发送失败.
const ErrCode40902 ErrCode = 40902 //	邮件处理发生未知异常
const ErrCode40903 ErrCode = 40903 //	邮件发送成功
const ErrCode40904 ErrCode = 40904 //	额度检查失败
const ErrCode40905 ErrCode = 40905 //	额度检查通过
const ErrCode40906 ErrCode = 40906 //	额度检查临时通过
const ErrCode40907 ErrCode = 40907 //	该API_USER对应的内容不需要进行模板匹配
const ErrCode40908 ErrCode = 40908 //	邮件内容和邮件模板匹配不通过
const ErrCode40909 ErrCode = 40909 //	邮件内容和邮件模板匹配通过
const ErrCode40910 ErrCode = 40910 //	邮件内容和邮件模板匹配临时通过
const ErrCode40911 ErrCode = 40911 //	邮件内容和邮件模板匹配时出现编码错误
const ErrCode41001 ErrCode = 41001 //	name不能为空串
const ErrCode41002 ErrCode = 41002 //	name的长度应该为1-250个字符
const ErrCode41003 ErrCode = 41003 //	name不符合域名规则
const ErrCode41004 ErrCode = 41004 //	newName不能为空串
const ErrCode41005 ErrCode = 41005 //	newName的长度应该为1-250个字符
const ErrCode41006 ErrCode = 41006 //	newName不符合域名规则
const ErrCode41007 ErrCode = 41007 //	type不能为空串
const ErrCode41008 ErrCode = 41008 //	type不符合规则
const ErrCode41009 ErrCode = 41009 //	verify不能为空串
const ErrCode41010 ErrCode = 41010 //	verify不符合规则
const ErrCode41011 ErrCode = 41011 //	verify解析错误
const ErrCode41012 ErrCode = 41012 //	用户创建域名不能超过5个
const ErrCode41013 ErrCode = 41013 //	name参数错误, 多个域名
const ErrCode41014 ErrCode = 41014 //	域名不存在
const ErrCode41015 ErrCode = 41015 //	domain创建失败
const ErrCode41016 ErrCode = 41016 //	domain修改失败
const ErrCode41101 ErrCode = 41101 //	emailType不能为空串
const ErrCode41102 ErrCode = 41102 //	emailType不符合规则
const ErrCode41103 ErrCode = 41103 //	cType不能为空串
const ErrCode41104 ErrCode = 41104 //	cType不符合规则
const ErrCode41105 ErrCode = 41105 //	domainName不能为空串
const ErrCode41106 ErrCode = 41106 //	domainName不符合规则
const ErrCode41107 ErrCode = 41107 //	domainName的长度应该为1-250个字符
const ErrCode41108 ErrCode = 41108 //	domainName所属的域名不存在
const ErrCode41109 ErrCode = 41109 //	用户信息不存在
const ErrCode41110 ErrCode = 41110 //	name不能为空串
const ErrCode41111 ErrCode = 41111 //	name不符合规则, name的长度为6-32的字符串, 只能含有(A-Z,a-z,0-9,_)
const ErrCode41112 ErrCode = 41112 //	apiUser不能超过10个
const ErrCode41113 ErrCode = 41113 //	open不能为空串
const ErrCode41114 ErrCode = 41114 //	open不符合规则
const ErrCode41115 ErrCode = 41115 //	click不能为空串
const ErrCode41116 ErrCode = 41116 //	click不符合规则
const ErrCode41117 ErrCode = 41117 //	unsubscribe不能为空串
const ErrCode41118 ErrCode = 41118 //	unsubscribe不符合规则
const ErrCode41119 ErrCode = 41119 //	apiUser创建失败
const ErrCode49901 ErrCode = 49901 //	url格式错误
const ErrCode49902 ErrCode = 49902 //	http请求执行异常
const ErrCode49903 ErrCode = 49903 //	http请求执行失败
const ErrCode49904 ErrCode = 49904 //	http请求执行成功
const ErrCode49905 ErrCode = 49905 //	http返回结果解析错误
const ErrCode49906 ErrCode = 49906 //	http其他错误
const ErrCode50000 ErrCode = 50000 //	接口频率受限(每个apiuser,每个接口、每分钟调用4000次，目前只限制投递回应)
const ErrCode50001 ErrCode = 50001 //	邮件发送失败.536 Frequency limited（每个apiuser每分钟调用总请求数不能超过4万次）
const ErrCode51001 ErrCode = 51001 //	sender不能为空
const ErrCode51002 ErrCode = 51002 //	sender长度不能超过250字符
const ErrCode51003 ErrCode = 51003 //	sender前缀不能包含@符号
const ErrCode51004 ErrCode = 51004 //	domain长度不能超过250字符
const ErrCode51005 ErrCode = 51005 //	指定删除数据不存在
const ErrCode51006 ErrCode = 51006 //	删除成功
const ErrCode51007 ErrCode = 51007 //	categoryName不能为空
const ErrCode51009 ErrCode = 51009 //	邮件日志查询未开通
const ErrCode501 ErrCode = 501     //	服务器异常
const ErrCode6001 ErrCode = 6001   //	你没有权限访问
