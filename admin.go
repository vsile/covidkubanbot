package main

import (
	"fmt"
	"log"
	"regexp"
	"net/http"
	"html/template"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

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

func corrLineBreaks(s string) string {
	re := regexp.MustCompile(`\n{3,}`)
	return re.ReplaceAllString(s, "\n\n")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)	//Добавляем в лог номер строки
	const dir string = "web"						//Путь к рабочей папке
	
    //Подключаемся к MongoDB
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	cc := session.DB("local").C("covidkuban_config")

	http.HandleFunc("/5e85cdfbbad7aeb709d6da92", func(w http.ResponseWriter, r *http.Request) {
		data := []rubric{}
		cc.Find(nil).All(&data)
	
		t, err := template.ParseFiles(dir+"/index.html")
		if err != nil {
			log.Println(err)
			return
		}
		t.Execute(w, data)
	})
	
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		id := r.PostFormValue("id")
		if !bson.IsObjectIdHex(id) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", id)
			return
		}
		index := r.PostFormValue("index")

		//contentEditable после первого разрыва строки удваивает реальное число разрывов
		answer := corrLineBreaks(r.PostFormValue("answer"))
		question := corrLineBreaks(r.PostFormValue("question"))
		err = cc.Update(bson.M{"categories.qas.id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{
			"categories.$.qas."+index+".question": question,
			"categories.$.qas."+index+".answer": answer,
		}})
		fmt.Fprint(w, err)
	})
	
	http.HandleFunc("/addQA", func(w http.ResponseWriter, r *http.Request) {
		categoryId := r.PostFormValue("categoryId")
		if !bson.IsObjectIdHex(categoryId) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", categoryId)
			return
		}
		newId := bson.NewObjectId()
		err = cc.Update(bson.M{"categories.id": bson.ObjectIdHex(categoryId)}, bson.M{"$push": bson.M{
			"categories.$.qas": bson.M{
				"$each": []bson.M{{
					"id": newId,
					"question": "",
					"answer": "",
				}},
				"$position": 0,
			},
		}})
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, newId.Hex())
	})
	
	http.HandleFunc("/deleteQA", func(w http.ResponseWriter, r *http.Request) {
		qaId := r.PostFormValue("qaId")
		if !bson.IsObjectIdHex(qaId) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", qaId)
			return
		}
		err = cc.Update(bson.M{"categories.qas.id": bson.ObjectIdHex(qaId)}, bson.M{"$pull": bson.M{
			"categories.$.qas": bson.M{
				"id": bson.ObjectIdHex(qaId),
			},
		}})
		fmt.Fprint(w, err)
	})
	
	http.HandleFunc("/addCategory", func(w http.ResponseWriter, r *http.Request) {
		rubricId := r.PostFormValue("rubricId")
		if !bson.IsObjectIdHex(rubricId) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", rubricId)
			return
		}
		newId := bson.NewObjectId()
		err = cc.UpdateId(bson.ObjectIdHex(rubricId), bson.M{"$push": bson.M{
			"categories": bson.M{
				"$each": []bson.M{{
					"id": newId,
					"title": r.PostFormValue("categoryName"),
					"qas": []qa{},
				}},
				"$position": 0,
			},
		}})
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, newId.Hex())
	})
	
	http.HandleFunc("/deleteCategory", func(w http.ResponseWriter, r *http.Request) {
		categoryId := r.PostFormValue("categoryId")
		if !bson.IsObjectIdHex(categoryId) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", categoryId)
			return
		}
		err = cc.Update(bson.M{"categories.id": bson.ObjectIdHex(categoryId)}, bson.M{"$pull": bson.M{
			"categories": bson.M{
				"id": bson.ObjectIdHex(categoryId),
			},
		}})
		fmt.Fprint(w, err)
	})
	
	http.HandleFunc("/editCategory", func(w http.ResponseWriter, r *http.Request) {
		categoryId := r.PostFormValue("categoryId")
		if !bson.IsObjectIdHex(categoryId) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", categoryId)
			return
		}
		err = cc.Update(bson.M{"categories.id": bson.ObjectIdHex(categoryId)}, bson.M{"$set": bson.M{
			"categories.$.name": r.PostFormValue("categoryName"),
		}})
		fmt.Fprint(w, err)
	})
	
	http.HandleFunc("/addRubric", func(w http.ResponseWriter, r *http.Request) {
		newId := bson.NewObjectId()
		err = cc.Insert(bson.M{
			"_id": newId,
			"title": r.PostFormValue("rubricName"),
			"categories": []category{},
		})
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, newId.Hex())
	})
	
	http.HandleFunc("/deleteRubric", func(w http.ResponseWriter, r *http.Request) {
		rubricId := r.PostFormValue("rubricId")
		if !bson.IsObjectIdHex(rubricId) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", rubricId)
			return
		}
		err = cc.RemoveId(bson.ObjectIdHex(rubricId))
		fmt.Fprint(w, err)
	})
	
	http.HandleFunc("/editRubric", func(w http.ResponseWriter, r *http.Request) {
		rubricId := r.PostFormValue("rubricId")
		if !bson.IsObjectIdHex(rubricId) {
			fmt.Fprintf(w, "Идентификатор %v не может быть преобразован в ObjectId", rubricId)
			return
		}
		err = cc.UpdateId(bson.ObjectIdHex(rubricId), bson.M{"$set": bson.M{
			"title": r.PostFormValue("rubricName"),
		}})
		fmt.Fprint(w, err)
	})
	
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {}) //Чтобы Chrome запускал сайт 1 раз

	http.Handle("/js/", http.FileServer(http.Dir(dir)))
	http.Handle("/css/", http.FileServer(http.Dir(dir)))
	http.Handle("/src/", http.FileServer(http.Dir(dir)))

	fmt.Println("Сервер успешно запущен!")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

