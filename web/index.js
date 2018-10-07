import './node_modules/normalize.css/normalize.css';
import './node_modules/@fortawesome/fontawesome-free/css/all.css';
import './node_modules/purecss/build/pure.css';
import './index.scss';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link, Switch } from 'react-router-dom';
import MainRoute from './components/main.js';
import NoMatch from './notfound.js';
import About from './about.js';
import SignIn from './users/sign_in.js';
import Oauth from './users/oauth.js';
import Home from './home/index.js';


ReactDOM.render((
  <Router>
    <div>
      <Switch>
        <MainRoute path='/' exact component={Home} />
        <MainRoute path='/about' component={About} />
        <Route path='/sign_in' component={SignIn} />
        <Route path='/oauth/callback' component={Oauth} />
        <Route component={NoMatch} />
      </Switch>
    </div>
  </Router>
), document.querySelector('#layout-container'));
