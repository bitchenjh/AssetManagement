package main

import (
	"github.com/fabric/core/chaincode/shim"
	"fmt"
	"github.com/fabric/protos/peer"
	"strconv"
)

type SimpleChaincode struct {

}


func (t *SimpleChaincode)Init(stub shim.ChaincodeStubInterface)peer.Response {
	_,args :=  stub.GetFunctionAndParameters()
	if len(args)!=2 {
		return shim.Error("实例化给定初始参数错误，必须为账户名称和初始金额")
	}
	v,err := strconv.Atoi(args[1])
	if err!=nil {
		return shim.Error("指定账户金额错误")
	}

	err := stub.PutState(args[0],[]byte(strconv.Itoa(v)))
	if err!=nil {
		return shim.Error("保存状态数据时发生错误")
	}

	return shim.Success(nil)

}

/*
invoke : set ,get
-c '{"Args",["set","jack","2000"]}'
-c '{"Args",["get","jack"]}'

*/

func (t *SimpleChaincode)Invove(stub shim.ChaincodeStubInterface)peer.Response {
	fun,args :=  stub.GetFunctionAndParameters()
	var result string
	var err error
	if fun == "set " {
		//修改账户余额代码
		result,err = set(stub,args)
	} else if fun == "get" {
		result,err = get(stub,args)
		//或者指定账户余额代码
	} else {
        return shim.Error("此次操作为非法操作")
	}

	if err != nil{
		return  shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}


func set(stub shim.ChaincodeStubInterface,args []string)(string,error){
	if len(args)!=2 {
		return "",fmt.Errorf("必须指定账户名和对应的金额")
	}
	//判断金额数据是否正确
    v,err := strconv.Atoi(args[1])
    if err!= nil{
    	return "",fmt.Errorf("指定金额错误：%s",err)
	}

	err = stub.PutState(args[0],[]byte(strconv.Itoa(v)))
	if err!=nil{
		return "",fmt.Errorf("保存状态数据发生错误：%s",err)
	}
	return "保存成功",nil
}

//根据指定的账户从账本中查询相应的数据
func get(stub shim.ChaincodeStubInterface,args []string)(string,error) {

	if len(args) != 1 {
		return "", fmt.Errorf("必须且只能指定账户名称")
	}

	result, err := stub.GetState(args[0])
	if err != nil {
       return "",fmt.Errorf("根据指定账户名称查询的状态数据时发生错误")
	}
	if result!=nil{
		return "",fmt.Errorf("根据指定账户没有查询到相应的状态数据")
	}
	return result,nil
}


func main(){
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("启动链码失败！")
	}
}
