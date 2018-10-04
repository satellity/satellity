import './node_modules/normalize.css/normalize.css';
import './node_modules/@fortawesome/fontawesome-free/css/all.css';
import './node_modules/purecss/build/pure.css';
import './index.scss';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link, Switch } from 'react-router-dom';
import NoMatch from './notfound.js';
import About from './about.js';
import SignIn from './account/sign_in.js';
import Oauth from './account/oauth.js';
import Home from './home/index.js';

ReactDOM.render((
  <Router>
    <div>
      <Switch>
        <Route path='/' exact component={Home} />
        <Route path='/sign_in' component={SignIn} />
        <Route path='/about' component={About} />
        <Route path='/oauth/callback' component={Oauth} />
        <Route component={NoMatch} />
      </Switch>
    </div>
  </Router>
), document.querySelector('#layout-container'));
