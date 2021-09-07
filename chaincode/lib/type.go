package lib
import(
	"time"
)
// Account 账户
type Account struct {
	UserId    string  `json:"userId"` //用户ID
	UserName  string  `json:"userName"`  //账号名
	Balance   float64 `json:"balance"`   //余额
}
// 定义资源及资源售卖
type Resource struct {
	UserId	     string  `json:"userId"`       //用户ID
	ResourceId   string  `json:"resourceId"`   //资源ID
	ResourceType string  `json:"resourceType"` //资源类型
	Price        float64 `json:"price"`        //资源单价
	ServiceTime  float64     `json:"serviceTime"`  //单词购买资源服务时长上限
}

//定义服务的记录
type ServiceRecord struct {
	ServiceId	string  `json:"serviceId"`  //服务ID
	BuyerId     string  `json:"buyerId"`    //用户ID
	SellerId    string  `json:"sellerId"`   //销售者ID
	Price       float64 `json:"price"`		//交易金额
	StartTime   time.Time `json:"startTime"`  //服务开始时间
	EndTime     time.Time `json:"endTime"`    //服务结束时间
}

const (
	AccountKey         = "account-key"
	SellingKey         = "selling-key"
	ServiceKey         = "service-key"
)
