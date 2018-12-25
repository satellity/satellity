import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import style from './style.scss';
import API from '../api/index.js';

class SiteWidget extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
  }

  componentDidMount() {
  }

  render() {
    // TODO
    let githubClientId = '71905afbd6e4541ad62b';
    if (process.env.NODE_ENV === 'development') {
      githubClientId = 'b9b78f343f3a5b0d7c99';
    }
    let signIn = '';
    if (!this.api.user.loggedIn()) {
      signIn = (
        <div className={style.sign_in}>
          <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${githubClientId}`}>Sign in with GitHub</a>
        </div>
      )
    }
    return (
      <div className={`${style.widget}`}>
        GoDiscourse is an open source community written in Go, get codes in <a href='https://github.com/godiscourse/godiscourse' target='blank' className='soft'>Github</a>.
          {signIn}
      </div>
    )
  }
}

export default SiteWidget;
