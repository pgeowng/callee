package render

import (
	"fmt"
	"html/template"
	"net/http"
)

var indexTmpl *template.Template
var indexTmplName = "index.tmpl"
var indexTmplPath = "./server/index.tmpl"

func init() {

	// https://github.com/meshhq/golang-html-template-tutorial/blob/master/main.go
	// thirdViewFuncMap := ThirdViewFormattingFuncMap()
	// thirdViewHTML := assets.MustAssetString("templates/third_view.html")
	// thirdViewTpl = template.Must(template.New("third_view").Funcs(thirdViewFuncMap).Parse(thirdViewHTML))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(200)
	fmt.Println(r)

	indexTmpl = template.Must(template.New(indexTmplName).ParseFiles(indexTmplPath))

	// render(w, r, indexTmpl)
	render(w, r, indexTmpl, indexTmplName, struct{}{})
}

// render(w, r, thirdViewTpl, "third_view", fullData)

func render(w http.ResponseWriter, r *http.Request, template *template.Template, name string, scope any) {
	fmt.Println(indexTmpl)
	err := template.Execute(w, scope)
	if err != nil {
		fmt.Println(err)
	}
}
