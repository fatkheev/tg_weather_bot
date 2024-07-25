package main

import (
    "fmt"
    "log"
    "os"
    "weather_bot/pkg/weather"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Ошибка при загрузке файла .env")
    }

    botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
    if botToken == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN не указан в файле .env")
    }

    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Авторизован под аккаунтом %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil {
            log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

            if update.Message.Text == "/start" {
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я ваш погодный бот.")
                button := tgbotapi.NewKeyboardButtonLocation("Показать погоду в текущем городе")
                row := tgbotapi.NewKeyboardButtonRow(button)
                keyboard := tgbotapi.NewReplyKeyboard(row)
                msg.ReplyMarkup = keyboard
                bot.Send(msg)
                continue
            }

            if update.Message.Location != nil {
                lat := update.Message.Location.Latitude
                lon := update.Message.Location.Longitude
                log.Printf("Получение погоды для координат: %f, %f", lat, lon)

                weatherData, err := weather.GetWeatherByCoordinates(lat, lon)
                if err != nil {
                    log.Printf("Ошибка при получении данных о погоде: %v", err)
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные о погоде.")
                    bot.Send(msg)
                    continue
                }

                response := fmt.Sprintf("Погода в %s:\nТемпература: %.2f°C\nВлажность: %d%%\nСкорость ветра: %.2f м/с",
                    weatherData.Name, weatherData.Main.Temp, weatherData.Main.Humidity, weatherData.Wind.Speed)
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
                bot.Send(msg)
            } else {
                city := update.Message.Text
                log.Printf("Получение погоды для города: %s", city)

                weatherData, err := weather.GetWeather(city)
                if err != nil {
                    log.Printf("Ошибка при получении данных о погоде: %v", err)
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные о погоде.")
                    bot.Send(msg)
                    continue
                }

                response := fmt.Sprintf("Погода в %s:\nТемпература: %.2f°C\nВлажность: %d%%\nСкорость ветра: %.2f м/с",
                    weatherData.Name, weatherData.Main.Temp, weatherData.Main.Humidity, weatherData.Wind.Speed)
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
                bot.Send(msg)
            }
        }
    }
}