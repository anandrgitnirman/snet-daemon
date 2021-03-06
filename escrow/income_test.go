package escrow

import (
	"fmt"
	price2 "github.com/singnet/snet-daemon/pricing"
	"github.com/singnet/snet-daemon/handler"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

type incomeValidatorMockType struct {
	err error
}

func (incomeValidator *incomeValidatorMockType) Validate(income *IncomeData) (err error) {
	return incomeValidator.err
}

type  MockPriceType struct{

}
func (priceType MockPriceType) GetPrice(GrpcContext *handler.GrpcStreamContext) (price *big.Int , err error) {
	return big.NewInt(0),nil
}

func TestIncomeValidate(t *testing.T) {
	one := big.NewInt(1)
	income := big.NewInt(0)

	pricing := &price2.PricingStrategy{}
	pricing.AddPricingTypes(&MockPriceType{})
	incomeValidator := NewIncomeValidator(pricing)
	price := big.NewInt(0)

	income.Sub(price, one)
	err := incomeValidator.Validate(&IncomeData{Income: income})
	msg := fmt.Sprintf("income %s does not equal to price %s", income, price)
	assert.Equal(t, NewPaymentError(Unauthenticated, msg), err)

	income.Set(price)
	err = incomeValidator.Validate(&IncomeData{Income: income})
	assert.Nil(t, err)

	income.Add(price, one)
	err = incomeValidator.Validate(&IncomeData{Income: income})
	msg = fmt.Sprintf("income %s does not equal to price %s", income, price)
	assert.Equal(t, NewPaymentError(Unauthenticated, msg), err)


}


type  MockPriceErrorType struct{

}
func (priceType MockPriceErrorType) GetPrice(GrpcContext *handler.GrpcStreamContext) (price *big.Int , err error) {
	return nil,fmt.Errorf("Error in Determining Price")
}

func TestIncomeValidateForPriceError(t *testing.T) {
	pricing := &price2.PricingStrategy{}
	pricing.AddPricingTypes(&MockPriceErrorType{})
	incomeValidator := NewIncomeValidator(pricing)
	err := incomeValidator.Validate(&IncomeData{Income: big.NewInt(0)})
	assert.Equal(t, err.Error(), "Error in Determining Price")

}
