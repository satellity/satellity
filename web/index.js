import './node_modules/normalize.css/normalize.css';
import './index.scss';
import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link, Switch } from "react-router-dom";
import About from './about.js';
import NoMatch from './notfound.js';
import Home from './home/index.js';

ReactDOM.render((
  <Router>
    <div>
      <Switch>
        <Route path='/' exact component={Home} />
        <Route path='/about' component={About} />
        <Route component={NoMatch} />
      </Switch>
    </div>
  </Router>
), document.querySelector('#layout-container'));
