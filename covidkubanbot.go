package main

import (
	"fmt"
	"log"
	"time"
	"strings"
	"net/http"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type subscriber struct {
	Id          bson.ObjectId   `_id`
    Chatid      int64
    User        *tgbotapi.User
    Step		int				//0 - обычный режим, 1 - режим ввода вопроса
    Regdate     time.Time
}

type rubric struct {
	Id			bson.ObjectId	`_id`
	Title		string
	Description	string
	Categories	[]category
}

type category struct {
	Id		bson.ObjectId
	Title	string
	Qas		[]qa
}

type qa struct {
	Id			bson.ObjectId
	Question	string
	Answer		string
}

type index struct {
	I	int
}
type mix []interface{}

func sendError(bot *tgbotapi.BotAPI, chatId int64, err error) {
	_, err = bot.Send(tgbotapi.NewMessage(chatId, err.Error()))
	if err != nil {
		log.Println(err)
	}
}

func createMainMenu(data []rubric) tgbotapi.ReplyKeyboardMarkup {
	keyboard := [][]tgbotapi.KeyboardButton{}
	row := []tgbotapi.KeyboardButton{}
	for i, v := range data {
		button := tgbotapi.NewKeyboardButton(v.Title)
		row = append(row, button)				//Добавляем кнопку в строку
		if i % 2 == 1 {
			keyboard = append(keyboard, row)	//Добавляем строку в клавиатуру
			row = []tgbotapi.KeyboardButton{}	//Очищаем массив-строку с кнопками
		}
	}
	//Добавляем в клавиатуру оставшиеся в массиве-строке кнопки
	if row != nil {
		keyboard = append(keyboard, row)
	}
	
	markup := tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: keyboard,
	}
	return markup
}

func getCategories(text string, data []rubric) []category {
	for _, v := range data {
		if v.Title == text {
			return v.Categories
		}
	}	
	return nil
}

//Если выбран раздел из меню второго уровня ищем только одну категорию
func getCategory(text string, data []rubric) category {
	for _, v := range data {
		for _, c := range v.Categories {
			if c.Title == text {
				return c
			}
		}
	}	
	return category{}
}

func getQuestions(categories []category) (string, []string, []string) {
	questions := ""
	n := 0
	qaIds := []string{} 	//Массив с идентификаторами вопросов-ответов
	titles := []string{}	//Массив с названиями категорий (меню второго уровня)
	for _, category := range categories {
		if category.Title != "" {
			titles = append(titles, category.Title)
		}
		for _, qa := range category.Qas {
			n++
			questions += fmt.Sprint(n)+". "+qa.Question+"\n\n"
			qaIds = append(qaIds, qa.Id.Hex())
		}
	}	
	return questions, qaIds, titles
}

func getQuestionsByCategory(qas []qa) (string, []string) {
	questions := ""
	qaIds := []string{} 	//Массив с идентификаторами вопросов-ответов
	for i, qa := range qas {
		questions += fmt.Sprint(i+1)+". "+qa.Question+"\n\n"
		qaIds = append(qaIds, qa.Id.Hex())
	}
	return questions, qaIds
}

//Создаем цифровую клавиатуру
func createNumericKeyboard(qaIds []string) tgbotapi.InlineKeyboardMarkup {
	keyboard := [][]tgbotapi.InlineKeyboardButton{}
	row := []tgbotapi.InlineKeyboardButton{}
	for i, qaId := range qaIds {
		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(i+1), qaId)
		row = append(row, button)					//Добавляем кнопку в строку
		if i % 6 == 5 {
			keyboard = append(keyboard, row)		//Добавляем строку в клавиатуру
			row = []tgbotapi.InlineKeyboardButton{}	//Очищаем массив строку с кнопками
		}
	}
	//Добавляем в клавиатуру оставшиеся в массиве-строке кнопки
	if row != nil {
		keyboard = append(keyboard, row)
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

//Создаем саб-меню
func createSubMenu(titles []string) tgbotapi.ReplyKeyboardMarkup {
	keyboard := [][]tgbotapi.KeyboardButton{}
	row := []tgbotapi.KeyboardButton{}
	for i, title := range titles {
		button := tgbotapi.NewKeyboardButton(title)
		row = append(row, button)				//Добавляем кнопку в строку
		if i % 2 == 1 {
			keyboard = append(keyboard, row)	//Добавляем строку в клавиатуру
			row = []tgbotapi.KeyboardButton{}	//Очищаем массив-строку с кнопками
		}
	}
	//Добавляем в клавиатуру оставшиеся в массиве-строке кнопки
	if row != nil {
		keyboard = append(keyboard, row)
	}
	//Добавляем кнопку назад
	button := tgbotapi.NewKeyboardButton("🔙 Назад")
	row = tgbotapi.NewKeyboardButtonRow(button)
	keyboard = append(keyboard, row)
	
	markup := tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: keyboard,
	}
	return markup
	
}

//Создаем кнопку "Назад" (при поиске по ключевым словам)
func createBackButton() tgbotapi.ReplyKeyboardMarkup {
	button := tgbotapi.NewKeyboardButton("🔙 Назад")
	row := tgbotapi.NewKeyboardButtonRow(button)
	return tgbotapi.NewReplyKeyboard(row)
}

func createYesNo(newQuestionId string) tgbotapi.InlineKeyboardMarkup {
	buttonYes := tgbotapi.NewInlineKeyboardButtonData("✅ Да", "yes_"+newQuestionId)
	buttonNo := tgbotapi.NewInlineKeyboardButtonData("❌ Нет", "no_")
	row := tgbotapi.NewInlineKeyboardRow(buttonYes, buttonNo)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}



func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)	//Добавляем в лог номер строки
	
	bot, err := tgbotapi.NewBotAPI("--hidden API KEY--")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("%s успешно авторизован!", bot.Self.UserName)
	
	
    updates := bot.ListenForWebhook("/"+bot.Token)
    go http.ListenAndServeTLS("0.0.0.0:88", "webhook_cert.pem", "webhook_pkey.pem", nil)
	
	/*bot.RemoveWebhook()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)*/
	
    //Привязываем mongodb к боту
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	cu := session.DB("local").C("covidkuban_users")
	cc := session.DB("local").C("covidkuban_config")
	cq := session.DB("local").C("covidkuban_newquest")

	for update := range updates {
		if update.Message != nil {
			upd := *update.Message
			
			//Получаем перечень разделов вопросов
			data := []rubric{}
			err := cc.Find(nil).All(&data)
			if err != nil {
				sendError(bot, upd.Chat.ID, err)
			}
			
			//Для начала проверяем на каком шаге находится пользователь
			user := subscriber{}
    		//Ищем пользователя по chatid
			err = cu.Find(bson.M{"chatid": upd.Chat.ID}).One(&user)
			//Если не находим, добавляем его в БД
			if err != nil {
				err = cu.Insert(bson.M{"_id": bson.NewObjectId(), "chatid": upd.Chat.ID, "user": upd.From, "regdate": time.Now()})
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
			}
			
			if user.Step == 1 {
				newId := bson.NewObjectId()
				cq.Insert(bson.M{"_id": newId, "question": upd.Text, "chatid": upd.Chat.ID, "from": user.User.FirstName+" "+user.User.LastName, "date": time.Now()})
				msg := tgbotapi.NewMessage(upd.Chat.ID, "Отправить Ваш вопрос на горячую линию?\n\n<b>От кого:</b> "+user.User.FirstName+" "+user.User.LastName+"\n<b>Дата и время:</b> "+time.Now().Format("02.01.2006 15:04")+"\n<b>Вопрос:</b> \""+upd.Text+"\"")
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = createYesNo(newId.Hex())
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				continue
			}

			
			if upd.Text == "/start" {
   				//Сбрасываем step пользователя
				cu.Update(bson.M{"chatid": upd.Chat.ID}, bson.M{"$set": bson.M{"step": 0}})
				
				msg := tgbotapi.NewMessage(upd.Chat.ID, "Здравствуйте!\n\nЗдесь Вы можете получить ответы на часто задаваемые вопросы, касающиеся ситуации с коронавирусом в Краснодарском крае.\n\nВыберите рубрику или введите ключевое слово.")
				msg.ReplyMarkup = createMainMenu(data)
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
			//Если пользователь выбрал один из разделов главного меню
			} else if categories := getCategories(upd.Text, data); categories != nil {
				questions, qaIds, titles := getQuestions(categories)
				if questions == "" {	//Такое может быть, если в базе данных нет вопросов
					msg := tgbotapi.NewMessage(upd.Chat.ID, "❕ Вопросов в рубрике пока нет.\n\nВы можете направить Ваш вопрос на горячую линию. Для этого отправьте команду /ask и введите Ваш вопрос.")
					_, err = bot.Send(msg)
					if err != nil {
						sendError(bot, upd.Chat.ID, err)
					}
					continue
				}
				
				msg := tgbotapi.NewMessage(upd.Chat.ID, "Выберите вопрос")
				//Добавляем саб-меню (если категории рубрик имеются)
				if len(titles) != 0 {
					msg.ReplyMarkup = createSubMenu(titles)
				}
				_, err := bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}

				msg = tgbotapi.NewMessage(upd.Chat.ID, questions)
				msg.ReplyMarkup = createNumericKeyboard(qaIds)
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
			
			//Если пользователь выбрал один из разделов меню второго уровня
			} else if cat := getCategory(upd.Text, data); cat.Id != "" {
				questions, qaIds := getQuestionsByCategory(cat.Qas)
				
				msg := tgbotapi.NewMessage(upd.Chat.ID, "Выберите вопрос")
				msg = tgbotapi.NewMessage(upd.Chat.ID, questions)
				msg.ReplyMarkup = createNumericKeyboard(qaIds)
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
							
			//Возвращаемся на меню верхнего уровня
			} else if upd.Text == "🔙 Назад" {
				msg := tgbotapi.NewMessage(upd.Chat.ID, "Выберите рубрику или введите ключевое слово.\n\nЕсли Вы не нашли Ваш вопрос, Вы можете направить его на горячую линию. Для этого отправьте команду /ask и введите Ваш вопрос.")
				msg.ReplyMarkup = createMainMenu(data)
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
			//Если пользователь вводит команду /ask
			} else if upd.Text == "/ask" {
				err := cu.Update(bson.M{"chatid": upd.Chat.ID}, bson.M{"$set": bson.M{"step": 1}})
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
					continue
				}
				msg := tgbotapi.NewMessage(upd.Chat.ID, "Пожалуста, введите Ваш вопрос")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
			//Поиск по ключевым словам
			} else {
				/*
					db.getCollection('covidkuban_config').aggregate([
						{ "$unwind": "$categories" },
						{ "$unwind": "$categories.qas" },
						{ "$match": { "$or": [
							{ "categories.qas.answer": /курортной/ },
							{ "categories.qas.question": /курортной/ }
						] } },
						{ "$project": { "id": "$categories.qas.id", "question": "$categories.qas.question" } }
					])
				*/
				qas := []qa{}
				err := cc.Pipe([]bson.M{
					{"$unwind": "$categories"},
					{"$unwind": "$categories.qas"},
					{"$match": bson.M{"$or": []bson.M{
						{"categories.qas.answer": bson.M{"$regex": bson.RegEx{Pattern: upd.Text, Options: "i"}}},
						{"categories.qas.question": bson.M{"$regex": bson.RegEx{Pattern: upd.Text, Options: "i"}}},
					}}},
					{"$project": bson.M{"id": "$categories.qas.id", "question": "$categories.qas.question"}},
				}).All(&qas)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
				if len(qas) == 0 {
					keyword := "❕ По заданному ключевому слову"
					if strings.Contains(upd.Text, " ") {
						keyword = "❕ По заданным ключевым словам"
					}
					msg := tgbotapi.NewMessage(upd.Chat.ID, keyword+" ничего не найдено. Попробуйте ввести другое ключевое слово.\n\nЕсли Вы не нашли Ваш вопрос, Вы можете направить его на горячую линию. Для этого отправьте команду /ask и введите Ваш вопрос.")
					_, err = bot.Send(msg)
					if err != nil {
						sendError(bot, upd.Chat.ID, err)
					}
					continue
				}
				
				questions, qaIds := getQuestionsByCategory(qas)
				
				msg := tgbotapi.NewMessage(upd.Chat.ID, "Выберите вопрос")
				msg.ReplyMarkup = createBackButton()	//Добавляем кнопку "Назад"
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}

				msg = tgbotapi.NewMessage(upd.Chat.ID, questions)
				msg.ReplyMarkup = createNumericKeyboard(qaIds)
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
			}
			
		} else if update.CallbackQuery != nil {
			upd := *update.CallbackQuery.Message
			callback := update.CallbackQuery.Data
			if !bson.IsObjectIdHex(callback) {
				//yes_5e88c5d2c3875d244c92823e или no_
				callbackData := strings.Split(callback, "_")
				//Исключаем возможные ошибки, связанные с нажатием "устаревшей" кнопки
				if len(callbackData) != 2 {
					continue
				}
				ps := "\n\n<i>Вопрос не отправлен</i>"
				closeText := "Чтобы направить другой вопрос, отправьте команду /ask и введите Ваш вопрос."
				if callbackData[0] == "yes" {
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "✅ Ваш вопрос успешно отправлен!"))
					ps = "\n\n<i>Вопрос отправлен</i>"
					closeText = "Благодарим за Ваш вопрос!"
					if bson.IsObjectIdHex(callbackData[1]) {
						cq.UpdateId(bson.ObjectIdHex(callbackData[1]), bson.M{"$set": bson.M{"sent": true}})
					}
					/*msg := tgbotapi.NewMessage(476580, "Пришел новый вопрос от пользователя:"+upd.Text[70:])
					_, err = bot.Send(msg)
					if err != nil {
						sendError(bot, upd.Chat.ID, err)
					}*/
				} else if callbackData[0] == "no" {
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ Вопрос не отправлен."))
				} else {
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Ошибка"))
				}
				//Сбрасываем step пользователя
				err := cu.Update(bson.M{"chatid": upd.Chat.ID}, bson.M{"$set": bson.M{"step": 0}})
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
				//Меняем текущее сообщение
				editmsg := tgbotapi.NewEditMessageText(upd.Chat.ID, upd.MessageID, upd.Text[70:]+ps)
				editmsg.ParseMode = "HTML"
				_, err = bot.Send(editmsg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
				//Получаем перечень разделов вопросов
				data := []rubric{}
				err = cc.Find(nil).Select(bson.M{"title": 1}).All(&data)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
				msg := tgbotapi.NewMessage(upd.Chat.ID, closeText)
				msg.ReplyMarkup = createMainMenu(data)
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
				
				continue
			}
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

			/*
				db.getCollection('covidkuban_config').aggregate([
					{ "$unwind": "$categories" },
					{ "$unwind": "$categories.qas" },
					{ "$match": { "categories.qas.id": ObjectId("5e85cf7d020c6189c1dd48aa") } },
					{ "$project": { "question": "$categories.qas.question", "answer": "$categories.qas.answer" } }
				])
			*/
			qa := qa{}
			err := cc.Pipe([]bson.M{
				{"$unwind": "$categories"},
				{"$unwind": "$categories.qas"},
				{"$match": bson.M{"categories.qas.id": bson.ObjectIdHex(callback)}},
				{"$project": bson.M{"question": "$categories.qas.question", "answer": "$categories.qas.answer"}},
			}).One(&qa)
			if err != nil {
				sendError(bot, upd.Chat.ID, err)
				continue
			}
			
			//Ищем индекс вопроса в массиве второго уровня
			/*db.getCollection('covidkuban_config').aggregate([
				{"$match": {"categories.qas.id": ObjectId("5e8639cabecfb896d7832280")}},
				{"$unwind": "$categories"},
				{
					"$project": {
						"i": { "$indexOfArray": [ "$categories.qas.id", ObjectId("5e8639cabecfb896d7832280") ] }
					}
				},
				{"$match": {"i": {"$gt":-1}}}
			])*/
			qaIndex := index{}
			cc.Pipe([]bson.M{
				{"$unwind": "$categories"},
				{"$project": bson.M{
					"i": bson.M{"$indexOfArray": mix{"$categories.qas.id", bson.ObjectIdHex(callback)}},
				}},
				{"$match": bson.M{"i": bson.M{"$gt":-1}}},
			}).One(&qaIndex)
			//Добавляем единицу к выбранному вопросу (для статистики популярности вопросов)
			cc.Update(bson.M{"categories.qas.id": bson.ObjectIdHex(callback)}, bson.M{"$inc": bson.M{"categories.$.qas."+fmt.Sprint(qaIndex.I)+".requests": 1}})
			
			//answerLength := len([]rune(qa.Answer))
			//numOfMsg := answerLength/4096+1
			answerArray := strings.Split(qa.Answer, "<cut>")	//Разбиваем длинный ответ на несколько сообщений
			for _, v := range answerArray {
				msg := tgbotapi.NewMessage(upd.Chat.ID, "<b>"+qa.Question+"</b>\n\n"+v)
				msg.ParseMode = "HTML"
			
				_, err = bot.Send(msg)
				if err != nil {
					sendError(bot, upd.Chat.ID, err)
				}
			}
		}
	}
}
