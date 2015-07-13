package main

import (
	_ "bytes"
	"fmt"
	"os"
	"github.com/maddyonline/pipe"
	_ "path"
	"strings"
)

type EmberResource struct {
	Name  string
	Properties []string
}

type EmberApp struct {
	Name string
	Resource *EmberResource
}

const ACTIVE_MODEL_ADAPTER = `import DS from 'ember-data';

export default DS.ActiveModelAdapter.extend({
	namespace: 'api'
});`

const FORM = `
<form {{action "save" on="submit"}}> <p>
<label>First Name:
{{input value=model.firstName}}
    </label>
  </p>
<p>
<label>Last Name:
      {{input value=model.lastName }}
    </label>
  </p>
 <p> <label>Email:
      {{input value=model.email}}
    </label>
  </p>
<p> <label>Twitter
      {{input value=model.twitter}}
    </label>
  </p>
<input type="submit" value="Save"/>
<button {{action "cancel"}}>Cancel</button> </form>
`

const INDEX = `
<h1>Friends Index</h1>
{{! The context here is the controller}} <h2>Total friends: {{model.length}}</h2>
<ul>
{{#each friend in model}}
<li>{{friend.id}} - {{friend.firstName}} {{friend.lastName}}</li>
{{/each}}
</ul>
`

const INDEX_ROUTE = `
import Ember from 'ember';
export default Ember.Route.extend({ model: function() {
return this.store.findAll('friend'); }
});
`

const ADD_NEW = `<h1>Adding New Friend</h1> {{partial "friends/form"}}`

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

func updateAdapter(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.Exec("ember", "g", "adapter", "application"),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(ACTIVE_MODEL_ADAPTER)),
			pipe.WriteFile("app/adapters/application.js", 0644),
		),
	)
	return p
}

func createForm(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.Exec("ember", "g", "template", fmt.Sprintf("%s/-form", app.Resource.Name)),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(FORM)),
			pipe.WriteFile("app/templates/friends/-form.hbs", 0644),
		),
	)
	return p
}

func createIndex(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.Exec("ember", "g", "route", fmt.Sprintf("%s/index", app.Resource.Name)),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(INDEX)),
			pipe.WriteFile("app/templates/friends/index.hbs", 0644),
		),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(INDEX_ROUTE)),
			pipe.WriteFile("app/routes/friends/index.js", 0644),
		),
	)
	return p
}

func createNew(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.Exec("ember", "g", "route", fmt.Sprintf("%s/new", app.Resource.Name)),
		pipe.Exec("ember", "g", "controller", fmt.Sprintf("%s/new", app.Resource.Name)),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(ADD_NEW)),
			pipe.WriteFile("app/templates/friends/new.hbs", 0644),
		),
	)
	return p
}

func createBasic(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		createIndex(app),
		createForm(app),
		//pipe.Exec("ember", "g", "route", fmt.Sprintf("%s/new", app.Resource.Name)),
		//pipe.Exec("ember", "g", "route", 
		//	fmt.Sprintf("%s/show", app.Resource.Name), 
		//	fmt.Sprintf("--path=:%s_id", strings.TrimSuffix(app.Resource.Name, "s")),),
		//pipe.Exec("ember", "g", "route", 
		//	fmt.Sprintf("%s/edit", app.Resource.Name), 
		//	fmt.Sprintf("--path=:%s_id/edit", strings.TrimSuffix(app.Resource.Name, "s")),),
		//pipe.Exec("ember", "g", "controller", fmt.Sprintf("%s/base", app.Resource.Name)),
		//pipe.Exec("ember", "g", "controller", fmt.Sprintf("%s/new", app.Resource.Name)),
		//pipe.Exec("ember", "g", "controller", fmt.Sprintf("%s/edit", app.Resource.Name)),
	)
	return p
}

func createEmberApp(app *EmberApp) pipe.Pipe {
	genResourceArgs := append([]string{"g", "resource", app.Resource.Name}, app.Resource.Properties...)
	p := pipe.Script(
		pipe.Exec("ember", "new", app.Name),
		pipe.ChDir(app.Name),
		pipe.Exec("ember", genResourceArgs...),
		updateAdapter(app),
		createBasic(app),
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