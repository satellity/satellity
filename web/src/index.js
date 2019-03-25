import '../node_modules/noty/src/noty.scss';
import '../node_modules/noty/src/themes/nest.scss';
import '../node_modules/normalize.css/normalize.css';
import '../node_modules/@fortawesome/fontawesome-free/css/all.css';
import './assets/css/h5bp.css';
import './index.scss';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link, Switch } from 'react-router-dom';
import showdown from 'showdown';
import MainLayout from './layouts/main.js';
import AdminRoute from './admin/admin.js';
import NoMatch from './notfound.js';
import Oauth from './users/oauth.js';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faComment, faComments, faEdit, faEye, faTrashAlt } from '@fortawesome/free-regular-svg-icons';
import { faEllipsisV, faPlus, faHome } from '@fortawesome/free-solid-svg-icons';
library.add(faComment, faComments, faEdit,
  faEye, faTrashAlt,
  faEllipsisV, faPlus, faHome);

showdown.setOption('customizedHeaderId', true);
showdown.setOption('simplifiedAutoLink', true);
showdown.setOption('strikethrough', true);
showdown.setOption('simpleLineBreaks', true);

ReactDOM.render((
  <Router>
    <div>
      <Switch>
        <Route path='/oauth/callback' component={Oauth} />
        <Route path='/admin' component={AdminRoute} />
        <Route path='/404' component={NoMatch} />
        <MainLayout />
      </Switch>
    </div>
  </Router>
), document.querySelector('#layout-container'));
