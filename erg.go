package main

import (
	_ "bytes"
	"fmt"
	"os"
	"github.com/maddyonline/pipe"
	_ "path"
	_ "strings"
)

type EmberResource struct {
	Name  string
	Properties []string
}

type EmberApp struct {
	Name string
	Resource *EmberResource
}

const script = `
ember new borrowers && \
ember g resource friends firstName:string lastName:string email:string twitter:string totalArticles:number 
ember g adapter application && \
ember g route friends/index && \
ember g route friends/new && \
ember g route friends/show --path=:friend_id && \
ember g route friends/edit --path=:friend_id/edit && \
ember g controller friends/base && \
ember g controller friends/new && \
ember g controller friends/edit && \
ember g template friends/-form 

`

func createEmberApp(app *EmberApp) pipe.Pipe {
	genResourceArgs := append([]string{"g", "resource", app.Resource.Name}, app.Resource.Properties...)
	p := pipe.Script(
		pipe.Exec("ember", "new", app.Name),
		pipe.ChDir(app.Name),
		pipe.Exec("ember", genResourceArgs...),
	)
	return p
}

func runScript(app *EmberApp) error {
	p := createEmberApp(app)
	s := pipe.NewState(os.Stdout, os.Stderr)
	s.Echo = true
	err := p(s)
	if err == nil {
		err = s.RunTasks()
	}
	return err
}


func main() {
	r := &EmberResource{
		"friends", 
		[]string{
			"firstName:string", 
			"lastName:string", 
			"email:string", 
			"twitter:string", 
			"totalArticles:number",},
		}
	app := &EmberApp{"borrowers", r}
	err := runScript(app)
	fmt.Println(err)
}