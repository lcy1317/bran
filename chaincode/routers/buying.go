package routers

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/lcy1317/bran/chaincode/lib"
	"github.com/lcy1317/bran/chaincode/utils"
	"strconv"
	"time"
)


func CreateBuying(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 3 {
		return shim.Error("参数个数不正确")
	}
	userIds := args[0]
	resourceIds := args[1]
	servicesTimes := args[2]
	// 参数的类型转换

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
		return shim.Error(fmt.Sprintf("购买者用户账户不存在！"))
	}

	var argsss []string
	argsss = append(argsss , string(resourceIds))
	resultss , err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingKey, argsss)
	if err != nil {
		return shim.Error(fmt.Sprintf("查询资源信息出错！"))
	}
	if resultss == nil {
		return shim.Error(fmt.Sprintf("购买的资源不存在！"))
	}
	//现在的情况是，资源存在，购买者的id也存在，我需要获得购买者的完整的account信息以及资源的信息

	var buyeruser lib.Account
	for _, v := range results {
		if v != nil {
			var sellingres lib.Account
			err := json.Unmarshal(v, &sellingres)
			if err != nil {
				return shim.Error(fmt.Sprintf("购买者账户反序列化出错: %s", err))
			}
			buyeruser = sellingres
		}
	}

	var restobuy lib.Resource
	for _, v := range resultss {
		if v != nil {
			var sellingres lib.Resource
			err := json.Unmarshal(v, &sellingres)
			if err != nil {
				return shim.Error(fmt.Sprintf("购买资源反序列化出错: %s", err))
			}
			restobuy = sellingres
		}
	}

    //查询销售者的账户了
	var args2 []string
	args2 = append(args2 , string(restobuy.UserId))
	results2 , err := utils.GetStateByPartialCompositeKeys2(stub, lib.AccountKey, args2)

	if err != nil {
		return shim.Error(fmt.Sprintf("查询用户信息出错！"))
	}
	if results2 == nil {
		return shim.Error(fmt.Sprintf("销售用户账户不存在！"))
	}
	//testout , _ := json.Marshal(restobuy.UserId)
	//return shim.Success([]byte(testout))
	var solderuser lib.Account
	for _, v := range results2 {
		if v != nil {
			var sellingres lib.Account
			err := json.Unmarshal(v, &sellingres)
			if err != nil {
				return shim.Error(fmt.Sprintf("反序列化出错: %s", err))
			}
			solderuser = sellingres
		}
	}
	if (buyeruser.UserId == solderuser.UserId){
		return shim.Error(fmt.Sprintf("不能购买自己的资源"))
	}

	if (tservicesTimes>restobuy.ServiceTime){
		return shim.Error(fmt.Sprintf("购买时长超出资源限制"))
	}
	if (buyeruser.Balance < restobuy.Price * tservicesTimes){
		return shim.Error(fmt.Sprintf("余额不足"))
	}
	//正式进入交易流程
	totuse := restobuy.Price * tservicesTimes
	buyeruser.Balance = buyeruser.Balance - totuse
	solderuser.Balance = solderuser.Balance +totuse
	//这交易就完成了，就要把整个交易记录写进账本
	nowtime := time.Now()
	ttime := nowtime.Unix()
	ttimeobj := time.Unix(ttime, 0)
	//fmt.Println(ttimeobj)
	etime := nowtime.Unix() + int64(tservicesTimes)
	etimeobj := time.Unix(etime, 0)
	//fmt.Println(etimeobj)
	serviceids := strconv.FormatInt(ttime, 10) + buyeruser.UserId
		//创建一个服务记录
	servicerecord := &lib.ServiceRecord{
		ServiceId: serviceids,
		BuyerId: buyeruser.UserId,
		SellerId: solderuser.UserId,
		Price: totuse,
		StartTime: ttimeobj,
		EndTime: etimeobj,
	}
	//写进账本，主键是一个复合主键
	if err := utils.WriteLedger(buyeruser,stub, lib.AccountKey, []string{buyeruser.UserId}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	if err := utils.WriteLedger(solderuser,stub, lib.AccountKey, []string{solderuser.UserId}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	if err := utils.WriteLedger(servicerecord, stub, lib.ServiceKey, []string{serviceids}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	realselling, err := json.Marshal(servicerecord)
	if err != nil {
		return shim.Error(fmt.Sprintf("资源购买的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(realselling)
	//return shim.Error(fmt.Sprintf("%s",string(realselling)))
}

func QueryServiceList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var sellingList []lib.ServiceRecord
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServiceKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var sellingres lib.ServiceRecord
			err := json.Unmarshal(v, &sellingres)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryServiceList-反序列化出错: %s", err))
			}
			sellingList = append(sellingList, sellingres)
		}
	}
	sellingListByte, err := json.Marshal(sellingList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryServiceList-序列化出错: %s", err))
	}
	return shim.Success(sellingListByte)
	//return shim.Error(fmt.Sprintf("%s",string(sellingListByte)))
}
