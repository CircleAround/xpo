package guestbook

import (
	"html/template"
	"net/http"
	"time"

	"fmt"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

// [START greeting_struct]
type Greeting struct {
	ID      string
	Author  string
	Content string
	Date    time.Time
}

// [END greeting_struct]

func init() {
	http.HandleFunc("/my", root)
	http.HandleFunc("/sign", sign)

	http.HandleFunc("/", Entry)
	http.HandleFunc("/loggedin", LoggedIn)
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

// [START func_root]
func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Ancestor queries, as shown here, are strongly consistent with the High
	// Replication Datastore. Queries that span entity groups are eventually
	// consistent. If we omitted the .Ancestor from this query there would be
	// a slight chance that Greeting that had just been written would not
	// show up in a query.
	// [START query]
	// q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Limit(10)
	q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
	// [END query]
	// [START getall]
	greetings := make([]Greeting, 0, 10)
	if _, err := q.GetAll(c, &greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// [END getall]
	if err := guestbookTemplate.Execute(w, greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// [END func_root]

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook</title>
  </head>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
			<pre>{{.Content}}</pre>
			<pre>{{.ID}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))

// [START func_sign]
func sign(w http.ResponseWriter, r *http.Request) {
	// [START new_context]
	c := appengine.NewContext(r)
	// [END new_context]
	id, ierr := uuid.NewRandom()
	if ierr != nil {
		http.Error(w, ierr.Error(), http.StatusInternalServerError)
		return
	}

	g := Greeting{
		ID:      id.String(),
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}
	// [START if_user]
	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}
	// We set the same parent key on every Greeting entity to ensure each Greeting
	// is in the same entity group. Queries across the single entity group
	// will be consistent. However, the write rate to a single entity group
	// should be limited to ~1/second.
	key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	// [END if_user]
}

// [END func_sign]

//////////////////////////////////////////////////////////

func Entry(w http.ResponseWriter, r *http.Request) {
	// コンテキスト生成
	c := appengine.NewContext(r)

	// ログイン用URL
	// ログイン後ユーザーを引数で渡したURLへリダイレクト
	_, err := user.LoginURL(c, "/loggined")

	// エラーハンドリング
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// テンプレートを表示
	tmpl := template.Must(template.New("enrty").Parse(entyTmpl))
	tmpl.Execute(w, nil)
}

// エントリーページのテンプレート
var entyTmpl = `
<!DOCTYPE html>
<html lang="ja">
<head>
	<meta charset="UTF-8">
	<title>ログイン前</title>
</head>
<body>
 <a href="/loggedin">ログインしてね</a>
</body>
</html>
`

// ログイン後のハンドラ
func LoggedIn(w http.ResponseWriter, r *http.Request) {
	// コンテキスト生成
	c := appengine.NewContext(r)

	// 現在ログインしているユーザーの情報を取得する
	u := user.Current(c)

	// 現在ログインしているユーザーがいない場合
	if u == nil {
		// ユーザーにサインインするように促すためのページのURLを返す
		// ログイン後ユーザーを引数で与えたURLにリダイレクトさせる
		url, _ := user.LoginURL(c, "/")
		fmt.Fprintf(w, `<a href="%s">ログイン用のサイトに移る</a>`, url)
		return
	}

	// ログアウト用のURL
	// ログインアウト後ユーザーを引数で渡したURLへリダイレクト
	_, err := user.LogoutURL(c, "/")

	// エラーハンドリング
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// テンプレートを表示
	tmpl := template.Must(template.New("loggedIn").Parse(loggedInTmpl))
	tmpl.Execute(w, nil)
}

// ログイン後のページのテンプレート
var loggedInTmpl = `
<!DOCTYPE html>
<html lang="ja">
<head>
	<meta charset="UTF-8">
	<title>ログイン後</title>
</head>
<body>
<a href="my">マイページ</a>
 <p>ログインできたよ</p>
 <a href="/">ログアウト</a>
</body>
</html>
`
