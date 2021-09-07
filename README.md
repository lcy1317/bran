# 合约设计

## 简单版本

房地产交易的示例没有给如何创建一个用户，他那个逻辑里没有创建用户所以写入账本的时候统一写入的是房地产信息。

功能需求：

- 创建一个账户（充值）`func CreateAccount(stub shim.ChaincodeStubInterface,args []string) pb.Response{}`
- 查询账户列表`func QueryAccountList(stub shim.ChaincodeStubInterface,args []string) pb.Response{}`
- 挂售资源`func CreateResource(stub shim.ChaincodeStubInterface,args []string) pb.Response{}`
- 查询在售资源`func QueryResourceList(stub shim.ChaincodeStubInterface,args []string) pb.Response{}`
- 下线一个资源`func DelResource(stub shim.ChaincodeStubInterface,args []string) pb.Response{}`
- 删除一个账户`func DelAccount(stub shim.ChaincodeStubInterface,args []string) pb.Response{}`
- 购买资源`func CreateBuying(stub shim.ChaincodeStubInterface, userid string, selleruserid string, buyservicetime int) pb.Response{}`
- 查询购买记录`func QueryServiceList(stub shim.ChaincodeStubInterface, args []string) pb.Response{}`

数据结构：

- 账户结构体：

```go
type Account struct {
	UserId    string  `json:"accountId"` //用户ID
	UserName  string  `json:"userName"`  //账号名
	Balance   float64 `json:"balance"`   //余额
}
```

- 售卖资源结构体：

```go
type Resource struct {
    UserId	     string  `json:"userId"`       //用户ID
	ResourceId   string  `json:"resourceId"`   //资源ID
    ResourceType string  `json:"resourceType"` //资源类型
	Price        float64 `json:"price"`        //资源单价
	ServiceTime  int     `json:"serviceTime"`  //单次购买资源服务时长上限
}
```

- 服务记录结构体

```go
type ServiceRecord struct {
    ServiceId	string  `json:"serviceId"`  //服务ID
	BuyerId     string  `json:"buyerId"`    //用户ID
    SellerId    string  `json:"sellerId"`   //销售者ID
	StartTime   float64 `json:"startTime"`  //服务开始时间
	EndTime     int     `json:"endTime"`    //服务结束时间
}
```

## go sdk/ rest api外界函数调用接口

|      函数名      |      功能      | 调用参数（均为string格式）                                   |                        备注                        |
| :--------------: | :------------: | ------------------------------------------------------------ | :------------------------------------------------: |
|  createAccount   |    创建账户    | userId, username, balance                                    |   不会检测账户是否存在，可以直接利用这个更新余额   |
| queryAccountList |  查询账户列表  | - 无参：查询全部账户数据 -有参：[]userid, 可以查询多个用户id的账户信息 | 返回全部账户，或者提供用户ID参数只返回那些用户信息 |
|  delAccountList  |    删除账户    | userid                                                       |         余额大于10账户不可删除，其余直接删         |
|  createSelling   |  创建挂售资源  | userid, resourceid, resourcetype, price, uplimittime         |       会检测资源是否存在，唯一标识为资源ID号       |
|  refreshSelling  |  更新挂售资源  | userid, resourceid, resourcetype, price, uplimittime         |      会更新已有资源，若资源没有则创建为新资源      |
| querySellingList |  查询在售资源  | - 无参：查询全部资源数据 - 有参：[]resourceid, 可以查询多个资源id的资源信息 |                  列出所有在售资源                  |
|  delSellingList  |  下线在售资源  | resourceid                                                   |                下线指定资源ID的资源                |
|   createBuying   |  建立一个交易  | buyeruserid, resourceid, buyservicetime                      |                      购买资源                      |
| queryServiceList | 查询已交易记录 | - 无参：查询全部交易数据 - 有参：[]serviceid, 可以查询多个服务交易id的资源信息 |        罗列所有交易记录，可以指定交易id查询        |

