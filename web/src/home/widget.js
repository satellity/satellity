import React, { Component } from 'react';
import Config from '../components/constants.js';
import style from './widget.scss';
import API from '../api/index.js';

class SiteWidget extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
  }

  render() {
    let signIn;
    if (!this.api.user.loggedIn()) {
      signIn = (
        <div className={style.sign_in}>
          <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${Config.GithubClientId}`}>Sign in with GitHub</a>
        </div>
      )
    }

    return (
      <div>
        <div className={style.widget}>
          <h2 className={style.site}>
            Go Discourse
          </h2>
          <ul className={style.features}>
            <li> 1. Open Source on <a href='https://github.com/godiscourse/godiscourse' target='blank' className='soft'>Github</a>. </li>
            <li> 2. Based on Golang, React and PostgreSQL. </li>
            <li> 3. Model tested. </li>
            <li> 4. Project <a href='https://github.com/godiscourse/godiscourse/projects/1' target='_blank'>Roadmap</a>. </li>
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
