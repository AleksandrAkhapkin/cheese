package service

import (
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"strings"
)

func (s *Service) finderCheeseInBill(billID int, items []types.Items) (bool, error) {

	names, err := s.p.GetNamesOfProduct()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetNamesOfProduct in finderCheeseInBill"))
		return false, err
	}
	for i := range items {
		for j := range names {
			if items[i].Name == names[j] || strings.Contains(items[i].Name, names[j]) {
				if err = s.p.MarkerCheesePosition(billID, items[i].Name); err != nil {
					logger.LogError(errors.Wrap(err, "err with MarkerCheesePosition in finderCheeseInBill"))
				}
				return true, nil
			}
		}
	}

	return false, nil
}
