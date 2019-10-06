import '../node_modules/noty/src/noty.scss';
import '../node_modules/noty/src/themes/nest.scss';
import '../node_modules/normalize.css/normalize.css';
import '../node_modules/@fortawesome/fontawesome-free/css/all.css';
import './assets/css/h5bp.css';
import './index.scss';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import showdown from 'showdown';
import Locale from './locale/index.js';
import MainLayout from './layouts/main.js';
import AdminRoute from './admin/admin.js';
import NoMatch from './sink.js';
import Oauth from './users/oauth.js';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faBookmark, faComment, faComments, faEdit, faEye, faTrashAlt, faHeart } from '@fortawesome/free-regular-svg-icons';
import { faChalkboard, faEllipsisV, faHome, faPlus, faUsersCog } from '@fortawesome/free-solid-svg-icons';
import { faMarkdown } from '@fortawesome/free-brands-svg-icons';
library.add(
  faBookmark, faComment, faComments,
  faEdit, faEye, faTrashAlt,
  faHeart,
  faChalkboard, faEllipsisV, faHome,
  faPlus, faUsersCog, faMarkdown
);

showdown.setOption('customizedHeaderId', true);
showdown.setOption('simplifiedAutoLink', true);
showdown.setOption('strikethrough', true);
showdown.setOption('simpleLineBreaks', true);

window.i18n = new Locale(navigator.language);

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
