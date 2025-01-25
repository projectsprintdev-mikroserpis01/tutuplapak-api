package service

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strconv"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/validator"
)

type purchaseService struct {
	repo      contracts.PurchaseRepository
	validator validator.ValidatorInterface
}

func NewPurchaseService(
	repo contracts.PurchaseRepository,
	validator validator.ValidatorInterface,
) contracts.PurchaseService {
	return &purchaseService{
		repo:      repo,
		validator: validator,
	}
}

func (s *purchaseService) Purchase(ctx context.Context, req dto.PurchaseRequest) (dto.PurchaseResponse, error) {
	var err error
	err = s.validateSenderContactDetail(req.SenderContactType, req.SenderContactDetail)
	if err != nil {
		return dto.PurchaseResponse{}, err
	}

	// Prepare purchase
	var purchasedItems []entity.PurchaseItem
	var totalPrice float64
	paymentDetails := make(map[int]dto.PaymentDetail)

	for _, item := range req.PurchasedItems {

		productId, err := strconv.Atoi(item.ProductID)
		if err != nil {
			return dto.PurchaseResponse{}, err
		}

		product, err := s.repo.GetProductById(ctx, productId)
		if err != nil {
			return dto.PurchaseResponse{}, err
		}

		// Check quantity
		if item.Qty > product.Quantity {
			return dto.PurchaseResponse{}, errors.New("quantity product less than purchased product")
		}

		// Add to purchased items
		purchasedItems = append(purchasedItems, entity.PurchaseItem{
			ProductID:        product.ID,
			Name:             product.Name,
			Category:         product.Category,
			Quantity:         item.Qty,
			Price:            product.Price,
			SKU:              product.SKU,
			FileID:           product.FileID,
			FileURL:          product.FileURL,
			FileThumbnailURL: product.FileThumbnailURL,
		})

		// Calculate total price
		totalPrice += float64(item.Qty) * product.Price

		seller, err := s.repo.GetSellerById(ctx, product.UserID)
		if err != nil {
			return dto.PurchaseResponse{}, err
		}

		paymentDetails[product.UserID] = dto.PaymentDetail{
			BankAccountName:   seller.BankAccountName,
			BankAccountHolder: seller.BankAccountHolder,
			BankAccountNumber: seller.BankAccountNumber,
			TotalPrice:        paymentDetails[product.UserID].TotalPrice + float64(item.Qty)*product.Price,
		}

	}
	purchaseId, err := s.repo.CreatePurchase(ctx, purchasedItems, req.SenderName, req.SenderContactType, req.SenderContactDetail)
	if err != nil {
		return dto.PurchaseResponse{}, err
	}

	// Flatten map values into a slice
	paymenDetailsSlice := make([]dto.PaymentDetail, 0, len(paymentDetails))
	for _, v := range paymentDetails {
		paymenDetailsSlice = append(paymenDetailsSlice, v)
	}

	return dto.PurchaseResponse{PurchaseID: strconv.FormatInt(purchaseId, 10), PurchasedItems: purchasedItems, TotalPrice: totalPrice, PaymentDetails: paymenDetailsSlice}, err
}

func (s *purchaseService) UploadPayment(ctx context.Context, req dto.UploadPaymentRequest, purchaseId string) error {

	id, err := strconv.Atoi(purchaseId)
	if err != nil {
		return err
	}
	purchase, err := s.repo.GetPurchaseById(ctx, id)
	if err != nil {
		return err
	}

	for _, purchasedItem := range purchase.PurchasedItems {
		// TODO: bulk decrease
		err := s.repo.DecreaseQuantity(ctx, purchasedItem.ProductID, purchasedItem.Quantity)
		if err != nil {
			return err
		}

	}
	return nil
}

func validatePhone(phone string) error {
	phoneRegex := `^\+\d{8,15}$`
	re := regexp.MustCompile(phoneRegex)
	if re.MatchString(phone) {
		return nil
	} else {
		return errors.New("invalid phone format")
	}
}

func (s *purchaseService) validateSenderContactDetail(contactType, contactDetail string) error {
	if contactType == "email" {
		_, err := mail.ParseAddress(contactDetail)
		return err
	} else if contactType == "phone" {
		return validatePhone(contactDetail)
	}
	return nil
}
