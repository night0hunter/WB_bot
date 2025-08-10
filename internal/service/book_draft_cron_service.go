package service

import (
	"context"
	"fmt"
	"os"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func (s *Service) BookDraftCron(ctx context.Context) (int64, error) {
	bookings, err := s.Repository.SelectBookings(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "Repository.SelectBookings")
	}

	for _, booking := range bookings {
		resp, err := s.Adapter.GetWarehouseGoodsV2(ctx, dto.GetWarehouseGoodsV2Request{
			DraftID:     booking.DraftID,
			WarehouseID: booking.Warehouse,
		}, os.Getenv("GETWHGOODSV2_URL"))
		if err != nil {
			return booking.ChatID, errors.Wrap(err, "Adapter.GetWarehouseGoodsV2")
		}

		ok := true
		for _, item := range resp.Items {
			if item.HasError == true {
				ok = false
				break
			}

			if booking.SupplyType == enum.Box {
				if item.CanMix != true {
					ok = false
					break
				}

				continue
			}

			if booking.SupplyType == enum.Monopallet {
				if item.CanMonopallet != true {
					ok = false
					break
				}
			}
		}

		var boxTypeMask = map[enum.SupplyType]int{
			enum.Box:        4,
			enum.Monopallet: 32,
		}

		if ok {
			resp, err := s.Adapter.Create(ctx, dto.GetCreateRequest{
				BoxTypeMask: boxTypeMask[booking.SupplyType],
				DraftID:     booking.DraftID,
				WarehouseID: booking.Warehouse,
			}, os.Getenv("CREATE_URL"))
			if err != nil {
				return booking.ChatID, errors.Wrap(err, "Adapter.Create")
			}

			if len(resp.IDs) > 1 {
				spew.Dump(resp)

				return booking.ChatID, &myError.MyError{
					ErrType: myError.PreorderIDLenError,
					Message: "preorderID: len != 1",
				}
			}

			err = s.Repository.UpdatePreorderID(ctx, booking.BookingID, resp.IDs[0].ID)
			if err != nil {
				return booking.ChatID, errors.Wrap(err, "Repository.UpdatePreorderID")
			}

			fmt.Println("Draft created successfully")

			// spew.Dump(resp)
		}
	}

	return 0, nil
}
