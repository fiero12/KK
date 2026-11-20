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

	// currencies := storage.NewStorage[string](5)
	// currencies.Add("BTC-RUB")
	// currencies.Add("ETH-RUB")
	// currencies.Add("USD-RUB")
	// currencies.Add("USD-RUB")

	// lengthOfCurrencies := currencies.Length()
	// logger.Info("Работа с валютами:",
	//  slog.Int("Количество уникальных комбинаций", lengthOfCurrencies))

	// traderStore := storage.NewStore[uuid.UUID, float64](10)
	// traderID := uuid.NewV7()

	// traderStore.Put(traderID, 50000.0004)
	// logger.Info("Новый трейдер добавлен в хранилище:",
	//  slog.String("UUID", traderID.String()),
	// slog.Float64("Начальный баланс", 50000.0004))


	// if balance, ok := traderStore.Get(traderID); ok{
	// 	logger.Info("поиск в хранилище:",
	//  	slog.String("найден трейдер с UUID", traderID.String()),
	// 	slog.Float64("его текущий баланс", balance))
	// }else{
	// 	logger.Error("элемент не найден в хранилище:",
	//  	slog.String("проверьте UUID", traderID.String()),)
	// }

	orderBook := storage.NewStore[uuid.UUID, *enum.Order](10)

	uuidOrder1 := uuid.NewV7()
	uuidOrder2 := uuid.NewV7()

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
						CreatedAt: time.Now().UTC(),
						UpdatedAt: time.Now().UTC(),
							}

	enum.LogOrderState(orderBuy, logger)

	orderSell := 	enum.Order{
						ID: uuidOrder2,
						TraderID: uuid.NewV4(),
						Side: enum.SideSell,
						Type: enum.OrderTypeLimit,
						Symbo: "USD",
						Price: 1000.00,
						Quantity: 2.00,
						Filled: 0.00,
						Status: enum.OrderStatusNew,
						CreatedAt: time.Now().UTC(),
						UpdatedAt: time.Now().UTC(),
						}

	enum.LogOrderState(orderSell, logger)
	fmt.Scanln()
	orderBook.Put(uuidOrder1, &orderBuy)
	orderBook.Put(uuidOrder2, &orderSell)

	fmt.Scanln()

	orderBook.Update(uuidOrder1, modifiers.ModifierForOrder(orderSell.Quantity), logger)
	orderBook.Update(uuidOrder2, modifiers.ModifierForOrder(orderBuy.Quantity), logger)

	logger.Info("Зарегистрирована сделка",
			slog.String("Type", orderBuy.Type.String()),
			slog.Float64("Value:", orderBuy.Price),
			slog.Float64("Quantity", orderBuy.Quantity),
			slog.String("Symbo", orderBuy.Symbo),
			slog.String("TraderID", orderBuy.TraderID.String()))
								

	fmt.Scanln()

}