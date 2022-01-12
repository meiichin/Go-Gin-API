package transaction

import (
	"goginapi/campaign"
	"goginapi/payment"
	"strconv"
)

type Service interface {
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProssesPayment(input TransactionNotificationInput) error
}

type service struct {
	repository         Repository
	paymentService     payment.Service
	campaignRepository campaign.Repository
}

func NewService(repository Repository, paymentService payment.Service, campaignRepository campaign.Repository) *service {
	return &service{repository, paymentService, campaignRepository}
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetTokenURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}
	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *service) ProssesPayment(input TransactionNotificationInput) error {
	transactionID, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.GetByID(transactionID)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_cart" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expired" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updateTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updateTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updateTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updateTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
