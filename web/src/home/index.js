import Style from './index.scss';
import React, {Component} from 'react';
import {Redirect} from 'react-router-dom';
import API from '../api/index.js';

class Index extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
  }

  render() {
    if (this.api.user.loggedIn()) {
      return (
        <Redirect to={{pathname: "/dashboard"}} />
      )
    }

    return (
      <h1>
        Hello, Go Discourse
      </h1>
    )
  }
}

export default Index;
