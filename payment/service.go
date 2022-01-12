package payment

import (
	"goginapi/user"
	"log"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetTokenURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetTokenURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-ztuJiSus63s_bqVpKQE1JB6V"
	midclient.ClientKey = "SB-Mid-client-ieUAdIjXpBtyJijw"
	midclient.APIEnvType = midtrans.Sandbox

	var snapGateway midtrans.SnapGateway
	snapGateway = midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			FName: user.Name,
			Email: user.Email,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		log.Println("error:", err)
		return "", err
	}
	log.Println("GetToken:", snapTokenResp)

	return snapTokenResp.RedirectURL, nil
}
