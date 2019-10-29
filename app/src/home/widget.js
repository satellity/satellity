import style from './widget.module.scss';
import React, { Component } from 'react';
import {Link} from 'react-router-dom';
import Config from '../components/config.js';
import Button from '../components/button.js';
import API from '../api/index.js';

class SiteWidget extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
  }

  render() {
    const i18n = window.i18n;
    let action;
    if (this.api.user.loggedIn()) {
      action = <Button type='link' action='/topics/new' text={i18n.t('topic.new')} classes='button' />;;
    }

    return (
      <div className={style.widget}>
        <div className={style.section}>
          <h2 className={style.site}>
            {Config.Name}
          </h2>
          <ul className={style.features} dangerouslySetInnerHTML={{__html: i18n.t('aside.rules')}}>
          </ul>
        </div>
        <div className={style.section}>
          <Link to='/avatar'>{i18n.t('avatar.title')}</Link>
        </div>
        <div className={style.action}>
          {action}
        </div>
        <div className={style.copyright}>
          © 2019 MIT license
        </div>
      </div>
    )
  }
}

export default SiteWidget;
