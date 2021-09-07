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


func CreateAccount(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 3 {
		return shim.Error("参数个数不满足")
	}
	userIds := args[0]
	userNames := args[1]
	balances := args[2]
	var formattedbalances float64
	if val, err := strconv.ParseFloat(balances, 64); err != nil {
		return shim.Error(fmt.Sprintf("余额参数格式转换出错: %s", err))
	} else {
		formattedbalances = val
	}
	account := &lib.Account{
		UserId: userIds,
		UserName:  userNames,
		Balance:   formattedbalances,
	}
	if err := utils.WriteLedger(account, stub, lib.AccountKey, []string{userIds}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	realaccount, err := json.Marshal(account)
	if err != nil {
		return shim.Error(fmt.Sprintf("账户序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(realaccount)
	//return shim.Error(fmt.Sprintf("%s",string(realaccount)))
}
func DelAccount(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 1 {
		return shim.Error("参数个数错误")
	}
	userIds := args[0]
	var argss []string
	argss = append(argss , string(userIds))
	results , err := utils.GetStateByPartialCompositeKeys2(stub, lib.AccountKey, argss)
	if err != nil {
		return shim.Error(fmt.Sprintf("查询账户信息出错"))
	}
	if results == nil {
		return shim.Error(fmt.Sprintf("账户不存在或已删除！"))
	}
	var ans []lib.Account
	for _, v := range results {
		if v != nil {
			var account lib.Account
			err := json.Unmarshal(v, &account)
			if err != nil {
				return shim.Error(fmt.Sprintf("DelAccountList-反序列化出错: %s", err))
			}
			ans = append(ans, account)
		}
	}
	if ans[0].Balance > 10 {
		return shim.Error(fmt.Sprintf("DelAccountList-余额大于10禁止删除"))
	}
	if err := utils.DelLedger(stub, lib.AccountKey, []string{userIds}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	return shim.Success([]byte("账户删除成功！"))
	//return shim.Error(fmt.Sprintf("账户删除成功"))
}
func QueryAccountList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountList []lib.Account
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var account lib.Account
			err := json.Unmarshal(v, &account)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryAccountList-反序列化出错: %s", err))
			}
			accountList = append(accountList, account)
		}
	}
	accountListByte, err := json.Marshal(accountList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAccountList-序列化出错: %s", err))
	}
	return shim.Success(accountListByte)
	//return shim.Error(fmt.Sprintf("%s",string(accountListByte)))
}
