package main

import (
	_ "bytes"
	"fmt"
	"os"
	"github.com/maddyonline/pipe"
	_ "path"
	"strings"
	"log"
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
<h2>Friends: {{model.length}}</h2>
<table>
  <thead>
<tr> <th>Name</th> <th></th>
    </tr>
  </thead>
<tbody>
    {{#each model as |friend|}}
<tr>
<td>{{link-to friend.firstName "friends.show" friend}}</td> <td><a href="#" {{action "delete" friend}}>Delete</a></td>
</tr>
{{/each}}
  </tbody>
</table>
`
const SHOW_PAGE = `
<ul>
<li>First Name: {{model.firstName}}</li>
<li>Last Name: {{model.lastName}}</li>
<li>Email: {{model.email}}</li>
<li>twitter: {{model.twitter}}</li>
<li>{{link-to "Edit info" "friends.edit" model}}</li>
<li><a href="#" {{action "delete" model}}>Delete</a></li>
</ul>

`

const INDEX_ROUTE = `
import Ember from 'ember';
export default Ember.Route.extend({ model: function() {
return this.store.findAll('friend'); }
});
`

const ADD_NEW = `<h1>Adding New Friend</h1> {{partial "friends/form"}}`
const ADD_NEW_ROUTE_JS = `import Ember from 'ember';
export default Ember.Route.extend({
  model: function() {
return this.store.createRecord('friend'); }
});`

const EDIT_PAGE = `<h1>Editing {{model.fullName}}</h1> {{partial 'friends/form'}}`
const BASE_CONTROLLER = `
import Ember from 'ember';
export default Ember.Controller.extend({ isValid: Ember.computed(
'model.email', 'model.firstName', 'model.lastName', 'twitter', function() {
return !Ember.isEmpty(this.get('model.email')) && !Ember.isEmpty(this.get('model.firstName')) && !Ember.isEmpty(this.get('model.lastName')) && !Ember.isEmpty(this.get('model.twitter'));
} ),
actions: {
save: function() {
if (this.get('isValid')) { var _this = this;
this.get('model').save().then(function(friend) { _this.transitionToRoute('friends.show', friend);
});
} else {
this.set('errorMessage', 'You have to fill all the fields'); }
return false; },
cancel: function() { return true;
} }
});
`

const FRIENDS_ROUTE = `
import Ember from 'ember';
export default Ember.Route.extend({ actions: {
delete: function(friend) { var _this = this;
friend.destroyRecord().then(function() { _this.transitionTo('friends.index');
}); }
} });
`

const ADD_NEW_CONTROLLER_JS = `
import FriendsBaseController from './base';
export default FriendsBaseController.extend({ actions: {
cancel: function() { this.transitionToRoute('friends.index'); return false;
} }
});
`

const EDIT_CONTROLLER = `
import FriendsBaseController from './base';
export default FriendsBaseController.extend({ actions: {
cancel: function() {
this.transitionToRoute('friends.show', this.get('model')); return false;
} }
});
`

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
		pipe.TeeLine(
			pipe.Read(strings.NewReader(ADD_NEW_ROUTE_JS)),
			pipe.WriteFile("app/routes/friends/new.js", 0644),
		),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(ADD_NEW_CONTROLLER_JS)),
			pipe.WriteFile("app/controllers/friends/new.js", 0644),
		),
	)
	return p
}


func createEdit(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.Exec("ember", "g", "route", 
			fmt.Sprintf("%s/edit", app.Resource.Name), 
			fmt.Sprintf("--path=:%s_id/edit", strings.TrimSuffix(app.Resource.Name, "s")),),
		pipe.Exec("ember", "g", "controller", fmt.Sprintf("%s/edit", app.Resource.Name)),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(EDIT_PAGE)),
			pipe.WriteFile("app/templates/friends/edit.hbs", 0644),
		),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(EDIT_CONTROLLER)),
			pipe.WriteFile("app/controllers/friends/edit.js", 0644),
		),
	)
	return p
}

func createShow(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.Exec("ember", "g", "route", 
			fmt.Sprintf("%s/show", app.Resource.Name), 
			fmt.Sprintf("--path=:%s_id", strings.TrimSuffix(app.Resource.Name, "s")),),
		pipe.Exec("ember", "g", "controller", fmt.Sprintf("%s/show", app.Resource.Name)),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(SHOW_PAGE)),
			pipe.WriteFile("app/templates/friends/show.hbs", 0644),
		),
	)
	return p
}

func createBaseController(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.Exec("ember", "g", "controller", fmt.Sprintf("%s/base", app.Resource.Name)),
		pipe.TeeLine(
			pipe.Read(strings.NewReader(BASE_CONTROLLER)),
			pipe.WriteFile("app/controllers/friends/base.js", 0644),
		),
	)
	return p
}

func createDelete(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		pipe.TeeLine(
			pipe.Read(strings.NewReader(FRIENDS_ROUTE)),
			pipe.WriteFile("app/routes/friends.js", 0644),
		),
	)
	return p
}

func createBasic(app *EmberApp) pipe.Pipe {
	p := pipe.Script(
		createIndex(app),
		createForm(app),
		createBaseController(app),
		createNew(app),
		createEdit(app),
		createShow(app),
		createDelete(app),
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func prepare(dest string) {
	yes, err := exists(dest)
	if err != nil {
		log.Fatal(err)
	}
	if yes {
		err = os.RemoveAll(dest)
		if err != nil {
			log.Fatal(err)
		}
	}
	os.Mkdir(dest, 0777)
}


func main() {
	prepare("new-app")
	os.Chdir("new-app")
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