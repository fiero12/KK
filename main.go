package main

import (
	"financialExchange/enum"
	"financialExchange/modifiers"
	"financialExchange/storage"
	"fmt"
	"log/slog"
	"os"
	"time"
	"uuid"
)

func main(){
	options := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, &options)
	logger := slog.New(handler)



	orderBook := storage.NewStore[uuid.UUID, enum.Order](10)

	tradeBook := storage.NewStore[uuid.UUID, enum.Trade](10)


	uuidOrder1 := uuid.NewV7()
	uuidOrder2 := uuid.NewV7()

	now := time.Now().UTC()

	orderBuy := 	enum.Order{
						ID: uuidOrder1,
						TraderID: uuid.NewV4(),
						Side: enum.SideBuy,
						Type: enum.OrderTypeLimit,
						Symbo: "USD",
						Price: 1000.00,
						Quantity: 2.00,
						Filled: 0.00,
						Status: enum.OrderStatusNew,
						CreatedAt: now,
						UpdatedAt: now,
							}

	enum.LogOrderState(orderBuy, logger)

	orderSell := 	enum.Order{
						ID: uuidOrder2,
						TraderID: uuid.NewV4(),
						Side: enum.SideSell,
						Type: enum.OrderTypeLimit,
						Symbo: "USD",
						Price: 1000.00,
						Quantity: 1.00,
						Filled: 0.00,
						Status: enum.OrderStatusNew,
						CreatedAt: now,
						UpdatedAt: now,
						}

	fmt.Scanln()

	logger.Info("Заявки созданы", "buy_order", orderBuy, "sell_order", orderSell)
	orderBook.Put(uuidOrder1, orderBuy)
	orderBook.Put(uuidOrder2, orderSell)

	fmt.Scanln()
}

	matchOrders(orderBook, orderBuy.ID, orderSell.ID, logger)

	// orderBook.Update(uuidOrder1, modifiers.ModifierForOrder(orderSell.Quantity), logger)
	// orderBook.Update(uuidOrder2, modifiers.ModifierForOrder(orderBuy.Quantity), logger)

	fmt.Scanln()

func matchOrders(storage *storage.Store[uuid.UUID, enum.order], buyID, sellID, logger *slog.logger){
	//достаем сделки из хранилища и потверждение успеха находки
	buyOrder, buyOk := storage.Get(buyID)// покупка и потверждение её находки
	sellOrder, sellOk := storage.Get(sellID)// продажа и потверждение её находки

	if !buyOk || !sellOk{// проверка на существование 
		logger.Error("Возникли ошибки с поиском заказов в хранилище")//логируем ошибку
		return
		// выходим из функции
	}

	if buyOrder.Symbo != sellOrder.Symbo{// проверка на совпадение валюты
		logger.Warn("Сделка не может быть осуществлена в связи с несовпадением валют")//логируем предупреждение
		return
		// выходим из функции
	}

	if buyOrder.Status == enum.OrderStatusFilled || sellOrder.Status == enum.OrderStatusFilled{
		// проверка статуса
		logger.Warn("Заказы уже завершены")//логируем предупреждение
		return
		// выходим из функции
	}
		
	if buyOrder.Price < sellOrder.Price{// проверка цен
		logger.Info("Сделка не может быть совершена из-за несовпадения цен")//логируем инфо
		return
		// выходим из функции
	}
		
	quantityTrade := min(buyOrder.Remaining(), sellOrder.Remaining())
	// вычисляем объем сделки

	updatedBuyOrder, ok:= storage.Update(buyID, modifiers.ModifierForOrder(quantityTrade))
	// обновление сделки на покупку

	if !ok{// обновление провалилось
		logger.Error("Произошла ошибка при обновлении buyOrder")// логируем ошибку
	}

	updatedSellOrder, ok:= storage.Update(sellID, modifiers.ModifierForOrder(quantityTrade))
	// обновление сделки на продажу
		if !ok{// обновление провалилось
		logger.Error("Произошла ошибка при обновлении sellOrder")// логируем ошибку
	}

	logger.Info("Заявки обновлены","buy_order", updatedBuyOrder, "sell_order", updatedSellOrder)
	// логируем успешно обновление сделок
	logger.Info("Сделка успешно проведена",
				slog.Float64("quantity", quantityTrade),
				slog.Any("buy_order", updatedBuyOrder),
				slog.Any("sell_order", updatedSellOrder),
			)// логируем изменения вроде
			// объем сделки
			// обновленные сделки
}
