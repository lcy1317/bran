package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"testing"
)

func initTest(t *testing.T) *shim.MockStub {
	scc := new(BlockChainRealEstate)
	stub := shim.NewMockStub("ex01", scc)
	checkInit(t, stub, [][]byte{[]byte("init")})
	return stub
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) pb.Response {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
	return res
}

// 测试链码初始化
func TestBlockChainRealEstate_Init(t *testing.T) {
	initTest(t)
}

func Test_QuerySellingList(t *testing.T) {
	stub := initTest(t)
	fmt.Println(fmt.Sprintf("1、测试新增资源\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("createSelling"),
			[]byte("00000001"),
			[]byte("202109031"),
			[]byte("2.4gHzWifi"),
			[]byte("9"),
			[]byte("100"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("2、测试新增资源\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("createSelling"),
			[]byte("00000002"),
			[]byte("202109032"),
			[]byte("5gHzWifi"),
			[]byte("99"),
			[]byte("100"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("3、测试新增资源\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("createSelling"),
			[]byte("00000002"),
			[]byte("202109033"),
			[]byte("HotPoint"),
			[]byte("9.5"),
			[]byte("100"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("4、测试获取所有数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("querySellingList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("5、测试新加入销售特定id的数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("querySellingList"),
			[]byte("202109032"),
			[]byte("202109031"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("6、测试新增账户\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("createAccount"),
			[]byte("00000006"),
			[]byte("wudi"),
			[]byte("1000"),
		}).Payload)))

	fmt.Println(fmt.Sprintf("6、测试新增资源\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("refreshSelling"),
			[]byte("00000006"),
			[]byte("202109032"),
			[]byte("HotPoint"),
			[]byte("23"),
			[]byte("100"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("7、测试获取所有数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("querySellingList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("8、测试账户删除\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("delSellingList"),
			[]byte("202109032"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("9、测试获取所有数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("querySellingList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("10、测试全部账户\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("11、测试资源购买\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("createBuying"),
			[]byte("00000003"),
			[]byte("202109031"),
			[]byte("10"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("12、测试全部账户\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("13、测试已交易服务\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryServiceList"),
		}).Payload)))
}

