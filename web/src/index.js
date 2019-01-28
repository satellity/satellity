import '../node_modules/normalize.css/normalize.css';
import '../node_modules/@fortawesome/fontawesome-free/css/all.css';
import './assets/css/h5bp.css';
import './index.scss';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link, Switch } from 'react-router-dom';
import showdown from 'showdown';
import MainRoute from './layouts/main.js';
import AdminRoute from './admin/admin.js';
import NoMatch from './notfound.js';
import Oauth from './users/oauth.js';
import Home from './home/index.js';
import UserEdit from './users/edit.js';
import UserTopics from './users/topics.js';
import TopicNew from './topics/new.js';
import TopicShow from './topics/show.js';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faEdit, faTrashAlt, faComment } from '@fortawesome/free-regular-svg-icons';
import { faPlus, faEllipsisV} from '@fortawesome/free-solid-svg-icons';
library.add(faEdit, faTrashAlt, faComment, faPlus, faEllipsisV);

showdown.setOption('customizedHeaderId', true);
showdown.setOption('simplifiedAutoLink', true);
showdown.setOption('strikethrough', true);
showdown.setOption('simpleLineBreaks', true);

ReactDOM.render((
  <Router>
    <div>
      <Switch>
        <MainRoute exact path='/' component={Home} />
        <MainRoute exact path='/user/edit' component={UserEdit} />
        <MainRoute path='/users/:id/topics' component={UserTopics} />
        <MainRoute exact path='/topics/new' component={TopicNew} />
        <MainRoute path='/topics/:id/edit' component={TopicNew} />
        <MainRoute path='/topics/:id' component={TopicShow} />
        <Route path='/oauth/callback' component={Oauth} />
        <Route path='/admin' component={AdminRoute} />
        <Route component={NoMatch} />
      </Switch>
    </div>
  </Router>
), document.querySelector('#layout-container'));
