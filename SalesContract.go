package main
import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"time"
)

type CoffeeMachine struct{
	CoffeeMachineID         string `json:"coffeeMachineID"`
	InstallationCompanyName string `json:"installationCompanyName"`
	InstallationAddress     string `json:"installationAddress"`
}
type FHCLContract struct{
	Contract struct {
		ContractID           string       `json:"contractID"`
		ContractDate         time.Time `json:"contractDate"`
		SalesManager         string    `json:"salesManager"`
		FHCLZonalManagerID   string    `json:"FHCLZonalManagerID"`
		FHCLZonalManagerName string    `json:"FHCLZonalManagerName"`
		ValidityInYears      string    `json:"validityInYears"`
		Customer             struct {
			CustomerName           string `json:"customerName"`
			CustomerID             string `json:"customerID"`
			CustomerAddress        string `json:"customerAddress"`
			CustomerRepresentative string `json:"customerRepresentative"`
		} `json:"Customer"`

		CoffeeMachines 		[]CoffeeMachine  `json:"CoffeeMachine"`
		CommercialTerms struct {
			CostPerMonth              string `json:"costPerMonth"`
			CoffeeCostPerKg           string `json:"coffeeCostPerKg"`
			TeaCostPerKg              string `json:"teaCostPerKg"`
			CostPerCupForNonFHCLBeans string `json:"costPerCupForNonFHCLBeans"`
			YearEndUplift             string `json:"yearEndUplift"`
		} `json:"CommercialTerms"`
		OtherTerms struct {
			MUCBillSettlementTime string `json:"MUCBillSettlementTime"`
			InterestRate          string `json:"interestRate"`
			DisputeRaiseTime      string `json:"disputeRaiseTime"`
			GracePeriod           string `json:"gracePeriod"`
			NoticePeriod          string `json:"noticePeriod"`
		} `json:"OtherTerms"`
	} `json:"Contract"`
}
type SalesContract struct {
}

func setSampleJSONValues(contract *FHCLContract) {
	contract.Contract.ContractID="1"
	contract.Contract.ContractDate=time.Now()
	contract.Contract.SalesManager="Luke Skywalker"
	contract.Contract.FHCLZonalManagerID="112233"
	contract.Contract.FHCLZonalManagerName  = "Obi Wan Kenobi"
	contract.Contract.ValidityInYears="5"
	contract.Contract.Customer.CustomerName="Death Star Inc"
	contract.Contract.Customer.CustomerID="99"
	contract.Contract.Customer.CustomerAddress="12/F1, Priya's Srinithya, 11th Street, Tansi Nagar, Velachery, Chennai - 600048"
	contract.Contract.Customer.CustomerRepresentative="Darth Vader"
	contract.Contract.CoffeeMachines=make([]CoffeeMachine,0)
	cm:=CoffeeMachine{CoffeeMachineID: "1",InstallationCompanyName: "xyz",InstallationAddress: "xxx",}
	contract.Contract.CoffeeMachines=append(contract.Contract.CoffeeMachines, cm)
	//contract.Contract.CoffeeMachine["1","Spectrum7","xyz"]
	/*contract.Contract.CoffeeMachine[0].CoffeeMachineID="1"
	contract.Contract.CoffeeMachine[0].InstallationCompanyName="Spectrum7"
	contract.Contract.CoffeeMachine[0].InstallationAddress="700 Parks Street, Mustafar"
*/	contract.Contract.CommercialTerms.CostPerMonth="5000"
	contract.Contract.CommercialTerms.CoffeeCostPerKg="50"
	contract.Contract.CommercialTerms.TeaCostPerKg="40"
	contract.Contract.CommercialTerms.CostPerCupForNonFHCLBeans="7"
	contract.Contract.CommercialTerms.YearEndUplift="10%"
	contract.Contract.OtherTerms.MUCBillSettlementTime="20"
	contract.Contract.OtherTerms.InterestRate="24%"
	contract.Contract.OtherTerms.DisputeRaiseTime="7"
	contract.Contract.OtherTerms.GracePeriod="30"
	contract.Contract.OtherTerms.NoticePeriod="30"

}
/* Implement the Init method from the chaincode interface */

func (t *SalesContract) Init(stub shim.ChaincodeStubInterface) peer.Response{
	args:=stub.GetStringArgs()
	var logger=shim.NewLogger("FHCLLogger")
	logger.Debugf("Number of arguments: %s", len(args))
	if len(args)==0{
		logger.Debugf("No of args is 0")
		var fhclContract FHCLContract
		setSampleJSONValues(&fhclContract)
		res,err:=json.Marshal(&fhclContract)
		if err!=nil {
			return shim.Error("Unable to marshal JSON")
		}

		myerr:=stub.PutState(fhclContract.Contract.ContractID,res)
		logger.Debugf("The contract ID is: %s", fhclContract.Contract.ContractID)
		logger.Debugf("The Sales Manager is: %s", fhclContract.Contract.SalesManager)
		if myerr!=nil{
			return shim.Error("Failed to create asset with the specified JSON")
		}
		//add code to marshall json
		//return shim.Error("Incorrect number of arguments")
		logger.Debugf("Number of args is 0, successfully setup block")
		return shim.Success(nil)
	} else if len(args)==2{
		err:=stub.PutState(args[0],[]byte(args[1]))
		if err!=nil{
			return shim.Error(fmt.Sprintf("Failed to create Asset: %s" , args[0]))
		}
	} else{
	return shim.Error("Incorrect number of arguments")
	}
	//err:=stub.PutState(args[0],[]byte(args[1]))
	return shim.Success(nil)

}
func (t *SalesContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	fn,args:=stub.GetFunctionAndParameters()
	var result string
	var err error

	if fn=="set"{
		result,err=set(stub,args)
	} else {
		result,err=get(stub,args)
	}
	if (err!=nil){
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(result))
}
func set(stub shim.ChaincodeStubInterface, args []string) (string,error){
	if (len (args) !=2){
		return "",fmt.Errorf("Wrong number of arguments supplied - Expected key and value")
	}
	err:=stub.PutState(args[0],[]byte(args[1]))
	if(err!=nil){
		return "", fmt.Errorf("Unable to set value for key: %s", args[0])
	}
	return args[1],nil
}
func get(stub shim.ChaincodeStubInterface, args []string) (string,error){
	if (len (args)!=1){
		return "", fmt.Errorf("Wrong number of arguments supplied - Expected only the key value")
	}
	response, err:=stub.GetState(args[0])
	if (err !=nil){
		return "", fmt.Errorf("Unable to get value for key: %s", args[0])
	}
	if (response==nil){
		return "", fmt.Errorf("Asset not found for %s:" , args[0])
	}
	return string(response),nil
}

func main() {
	if err:=shim.Start(new(SalesContract)); err!=nil {
		fmt.Printf("Error starting SalesContract chaincode: %s", err)
	}
}
