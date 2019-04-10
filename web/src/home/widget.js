import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import Config from '../components/config.js';
import style from './widget.scss';
import API from '../api/index.js';

class SiteWidget extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
  }

  render() {
    let action = <div className={style.newTopic}> <Link to='/topics/new'>{i18n.t('topic.new')}</Link> </div>;
    if (!this.api.user.loggedIn()) {
      action = (
        <div className={style.signIn}>
          <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${Config.GithubClientId}`}>{i18n.t('login.github')}</a>
        </div>
      )
    }

    return (
      <div className={style.widget}>
        <div className={style.section}>
          <h2 className={style.site}>
            Go Discourse
          </h2>
          <ul className={style.features} dangerouslySetInnerHTML={{__html: i18n.t('aside.rules')}}>
          </ul>
        </div>
        {action}
        <div className={style.copyright}>
          © 2019 MIT license
        </div>
      </div>
    )
  }
}

export default SiteWidget;
