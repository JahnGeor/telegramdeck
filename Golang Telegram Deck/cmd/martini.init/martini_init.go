package martini_init

import (
	"GolangTD/database"
	"GolangTD/session"
	"GolangTD/utils"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

const COOKIE_NAME = "User"

func InitMartini() *martini.ClassicMartini { //функция инициализатор работы http-handler (martini) (настройки Handler)
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  "../ui/templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
		IndentJSON: true,
	}))

	m.Use(martini.Static("../ui/static", martini.StaticOptions{Prefix: "static"})) //обработка url-запросов для статичных файлов
	m.Use(martini.Static("../data", martini.StaticOptions{Prefix: "data"}))
	return m
}

func CommonHandler(m *martini.ClassicMartini) {

	m.Get("/", indexHandler)      //обработчик запроса главной страницы
	m.Get("/posts", postsHandler) //обработчик
	m.Get("/edit", editHandler)
	m.Get("/delete", deleteHandler)           //обработчик
	m.Get("/post", postHandler)               //обработчик
	m.Get("/help", helpHandler)               //обработчик
	m.Get("/login", loginHandler)             //обработчик
	m.Get("/logout", logoutHandler)           //обработчик
	m.Get("/publish", publishHandler)         //обработчик
	m.Get("/editprofile", editProfileHandler) //
	m.Get("/error", errorHandler)
	m.Post("/login", postLoginHandler)
	m.Post("/register", postRegisterHandler)
	m.Post("/publish", postPublishHandler)
	m.Post("/edit", postEditHandler)
	m.Run()
}

//Post-обработчики
func postLoginHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	user, err := database.Authorisation(login, password)
	if err != nil {
		rnd.Redirect(errorCode("999"))
	} else {

		sessionID := session.SessionInMemory.Init(user.UserID)
		cookie := &http.Cookie{
			Name:    "SessionID",
			Value:   sessionID,
			Expires: time.Now().Add(5 * time.Minute),
		}
		http.SetCookie(w, cookie)
		fmt.Println("Значение ID: " + fmt.Sprint(session.SessionInMemory.Get(cookie.Value)))
		rnd.Redirect("/")
	}
}

func postRegisterHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")
	tokenID := r.FormValue("token")
	date := r.FormValue("date")
	fmt.Println()
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	u := database.NewUser(utils.GeneratedID(), login, firstname, lastname, date, tokenID, password, email)
	if u.Create() != nil {
		rnd.Redirect(errorCode("Пользователь с данной почтой или логином уже зарегистрирован!"))
	} else {
		rnd.Redirect("/login")
	}
}

func postPublishHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	data, err := database.FormData(r)
	fmt.Print("Hello")
	if err != nil {
		rnd.Redirect(errorCode(fmt.Sprint(err)))
	} else {
		title := r.FormValue("title")
		content := r.FormValue("content")
		date := r.FormValue("date")
		dateCreate := time.Now().String()
		userID := data.DataUser.UserID
		contentID := utils.GeneratedID()
		fmt.Println(createPhoto("posts", contentID, r))
		p := database.NewPost(title, content, dateCreate, date, contentID, userID)
		if err = p.Create(); err != nil {
			rnd.Redirect(errorCode(fmt.Sprint(err)))
		} else {
			rnd.Redirect("/posts")
		}

	}
}

func postEditHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	contentID := r.FormValue("contentID")
	date := r.FormValue("date")
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		os.Remove("../data/img/posts/" + contentID + ".jpg")
		dst, _ := os.Create("../data/img/" + "posts" + "/" + contentID + ".jpg")
		io.Copy(dst, file)
		defer dst.Close()
	}
	err = database.EditPost(title, content, date, contentID)
	rnd.Redirect("/posts")
}

//Get-обработчики
func indexHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	cookie, err := session.CheckCookie(r)
	if err != nil {
		rnd.Redirect("/login")
	} else {
		fmt.Print(cookie.Value)
		if session.SessionInMemory.Get(cookie.Value) == nil {
			rnd.Redirect("/login")
		}
		userID := session.SessionInMemory.Get(cookie.Value).UserID

		posts, err := database.FindPostsByUserID(userID)
		if err != nil {
			rnd.Redirect(errorCode("404"))
		}
		user := &database.User{}
		err = user.Get(userID)
		data := &database.Data{posts, user}
		for _, content := range data.DataPosts {
			content.ShowDeb()
		}
		data.DataUser.ShowDeb()
		rnd.HTML(200, "index", data)
	}

}

func postsHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	_, err := session.CheckCookie(r)
	if err != nil {
		rnd.Redirect("/login")
	} else {
		data, err := database.FormData(r)
		if err != nil {
			rnd.Redirect(errorCode(fmt.Sprint(err)))
		} else {
			rnd.HTML(200, "posts", data)
		}
	}

}

func postHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	_, err := session.CheckCookie(r)
	if err != nil {
		rnd.Redirect("/login")
	} else {
		contentID := r.URL.Query().Get("id")
		p := &database.Post{}
		p.FindPostByContentID(contentID)
		p.ShowDeb()
		rnd.HTML(200, "post", p)
	}
}

func helpHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	_, err := session.CheckCookie(r)
	if err != nil {
		rnd.Redirect("/login")
	} else {
		rnd.HTML(200, "help", nil)
	}
}

func loginHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	tr, err := template.ParseFiles("../ui/templates/login.html")
	if err != nil {
		rnd.Redirect(errorCode("404"))
	} else {
		tr.Execute(w, nil)
	}
}

func logoutHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "SessionID",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	rnd.Redirect("/login")
}

func publishHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	_, err := session.CheckCookie(r)
	if err != nil {
		rnd.Redirect("/login")
	} else {
		rnd.HTML(200, "publish", nil)
	}
}

func editProfileHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {

}

func editHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	cookie, err := session.CheckCookie(r)
	if err != nil {
		rnd.Redirect("/login")
	} else {
		contentID := r.URL.Query().Get("id")
		if database.PostDeleteCheck(session.SessionInMemory.Get(cookie.Value).UserID, contentID) == nil {
			p := &database.Post{}
			if err := p.FindPostByContentID(contentID); err != nil {
				rnd.Redirect(errorCode("Проблема с нахождением поста в базе."))
			} else {
				rnd.HTML(200, "edit", p)
			}
		} else {
			rnd.Redirect(errorCode("У вас не прав для совершения данной операции!"))
		}

	}
}

func deleteHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	cookie, err := session.CheckCookie(r)
	if err != nil {
		rnd.Redirect("/login")
	} else {
		userID := session.SessionInMemory.Get(cookie.Value).UserID
		id := r.URL.Query().Get("id")
		if err = database.PostDeleteCheck(userID, id); err != nil {
			rnd.Redirect(errorCode(fmt.Sprint(err)))
		} else {
			fileURL := "../data/img/posts/" + id + ".jpg"
			fmt.Print(fileURL)
			os.Remove(fileURL)
			err := database.DeletePost(id)
			if err != nil {
				rnd.Redirect(errorCode(fmt.Sprint(err)))
			}
			rnd.Redirect("/posts")
		}
	}

}

func errorHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// decryptCookie(r)
	rnd.HTML(200, "error", code)
}

func errorCode(code string) string {
	return "/error?code=" + code
}

func createPhoto(class string, id string, r *http.Request) error {
	file, _, err := r.FormFile("image")
	if err != nil {
		return err
	}
	defer file.Close()
	dst, err := os.Create("../data/img/" + class + "/" + id + ".jpg")
	if err != nil {
		return err
	}
	io.Copy(dst, file)
	defer dst.Close()
	return nil
}
