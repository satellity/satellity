import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import Config from './constants.js';
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
    let signIn = '';
    if (!this.api.user.loggedIn()) {
      signIn = (
        <div className={style.sign_in}>
          <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${Config.GithubClientId()}`}>Sign in with GitHub</a>
        </div>
      )
    }
    return (
      <div>
        <div className={style.widget}>
          <div className={style.name}>
            <FontAwesomeIcon icon={['far', 'comments']} />
            GoDiscourse
          </div>
          <ul className={style.features}>
            <li> 1. Open Source on <a href='https://github.com/godiscourse/godiscourse' target='blank' className='soft'>Github</a>. </li>
            <li> 2. Based on Golang, React and PostgreSQL. </li>
            <li> 3. Model tested. </li>
          </ul>
        </div>
        {signIn}
        <div className={style.copyright}>
          © 2019 MIT license
        </div>
      </div>
    )
  }
}

export default SiteWidget;
