package routers

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/lcy1317/bran/chaincode/lib"
	"github.com/lcy1317/bran/chaincode/utils"
	"strconv"
)


func CreateSelling(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 5 {
		return shim.Error("参数个数不满足")
	}
	userIds := args[0]
	resourceIds := args[1]
	resourceTypes := args[2]
	prices := args[3]
	servicesTimes := args[4]
	// 参数的类型转换
	var tprices float64
	if val, err := strconv.ParseFloat(prices, 64); err != nil {
		return shim.Error(fmt.Sprintf("价格参数格式转换出错: %s", err))
	} else {
		tprices = val
	}

	var tservicesTimes float64
	if val, err := strconv.ParseFloat(servicesTimes, 64); err != nil {
		return shim.Error(fmt.Sprintf("服务时长参数格式转换出错: %s", err))
	} else {
		tservicesTimes = val
	}

	var argss []string
	argss = append(argss , string(userIds))
	results , err := utils.GetStateByPartialCompositeKeys2(stub, lib.AccountKey, argss)

	if err != nil {
		return shim.Error(fmt.Sprintf("查询用户信息出错！"))
	}
	if results == nil {
		return shim.Error(fmt.Sprintf("该用户不存在！"))
	}

	var argsss []string
	argsss = append(argsss , string(resourceIds))
	resultss , err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingKey, argsss)

	if err != nil {
		return shim.Error(fmt.Sprintf("查询资源信息出错！"))
	}
	if resultss != nil {
		return shim.Error(fmt.Sprintf("该资源已存在！如需更新资源请使用refreshSelling"))
	}

	//创建一个资源的销售
	resourceselling := &lib.Resource{
		UserId: userIds,
		ResourceId: resourceIds,
		ResourceType: resourceTypes,
		Price: tprices,
		ServiceTime: tservicesTimes,
	}
	//写进账本，主键是一个复合主键
	if err := utils.WriteLedger(resourceselling, stub, lib.SellingKey, []string{resourceIds}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	realselling, err := json.Marshal(resourceselling)
	if err != nil {
		return shim.Error(fmt.Sprintf("资源序列化创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(realselling)
	//return shim.Error(fmt.Sprintf("%s",string(realselling)))
}
func RefreshSelling(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 5 {
		return shim.Error("参数个数不满足")
	}
	userIds := args[0]
	resourceIds := args[1]
	resourceTypes := args[2]
	prices := args[3]
	servicesTimes := args[4]
	// 参数的类型转换
	var tprices float64
	if val, err := strconv.ParseFloat(prices, 64); err != nil {
		return shim.Error(fmt.Sprintf("价格参数格式转换出错: %s", err))
	} else {
		tprices = val
	}

	var tservicesTimes float64
	if val, err := strconv.ParseFloat(servicesTimes, 64); err != nil {
		return shim.Error(fmt.Sprintf("服务时长参数格式转换出错: %s", err))
	} else {
		tservicesTimes = val
	}

	var argss []string
	argss = append(argss , string(userIds))
	results , err := utils.GetStateByPartialCompositeKeys2(stub, lib.AccountKey, argss)

	if err != nil {
		return shim.Error(fmt.Sprintf("查询用户信息出错！"))
	}
	if results == nil {
		return shim.Error(fmt.Sprintf("该用户不存在！"))
	}

	//创建一个资源的销售
	resourceselling := &lib.Resource{
		UserId: userIds,
		ResourceId: resourceIds,
		ResourceType: resourceTypes,
		Price: tprices,
		ServiceTime: tservicesTimes,
	}
	//写进账本，主键是一个复合主键
	if err := utils.WriteLedger(resourceselling, stub, lib.SellingKey, []string{resourceIds}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	realselling, err := json.Marshal(resourceselling)
	if err != nil {
		return shim.Error(fmt.Sprintf("资源序列化创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(realselling)
	//return shim.Error(fmt.Sprintf("%s",string(realselling)))
}
func DelSelling(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 1 {
		return shim.Error("参数个数错误")
	}
	resourceIds := args[0]
	var argss []string
	argss = append(argss , string(resourceIds))
	results , err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingKey, argss)

	if err != nil {
		return shim.Error(fmt.Sprintf("查询资源信息出错"))
	}
	if (len(results) < 1) {
		return shim.Error(fmt.Sprintf("该资源不存在或已删除！"))
	}
	//fmt.Println(len(results))
	var ans []lib.Resource
	for _, v := range results {
		if v != nil {
			var sellingres lib.Resource
			err := json.Unmarshal(v, &sellingres)
			if err != nil {
				return shim.Error(fmt.Sprintf("DelSellingList-反序列化出错: %s", err))
			}
			ans = append(ans, sellingres)
		}
	}

	if err := utils.DelLedger(stub, lib.SellingKey, []string{resourceIds}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	return shim.Success([]byte("资源删除成功！"))
	//return shim.Success([]byte(strconv.Itoa(len(results))))
	//return shim.Error(fmt.Sprintf("资源删除成功"))
}
func QuerySellingList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var sellingList []lib.Resource
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.SellingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var sellingres lib.Resource
			err := json.Unmarshal(v, &sellingres)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingList-反序列化出错: %s", err))
			}
			sellingList = append(sellingList, sellingres)
		}
	}
	sellingListByte, err := json.Marshal(sellingList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingList-序列化出错: %s", err))
	}
	return shim.Success(sellingListByte)
	//return shim.Error(fmt.Sprintf("%s",string(sellingListByte)))
}
