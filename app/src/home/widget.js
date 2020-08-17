import style from './widget.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import Button from '../components/button.js';
import API from '../api/index.js';

class SiteWidget extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
  }

  render() {
    const i18n = window.i18n;

    return (
      <div className={style.widget}>
        {
          this.api.user.loggedIn() && (
            <div className={style.new}>
              <Button type='link' action='/topics/new' text={i18n.t('topic.new')} classes='button' />
            </div>
          )
        }
        <div className={style.section}>
          <h2 className={style.title}>
            {i18n.t('general.welcome')}
          </h2>
          <ul className={style.rules} dangerouslySetInnerHTML={{__html: i18n.t('aside.rules')}}>
          </ul>
        </div>
        <div className={style.avatar}>
          <Link to='/products'>Person Creator Collection</Link>
        </div>
        <div className={style.copyright}>
          © 2019 - Now MIT license
        </div>
      </div>
    )
  }
}

export default SiteWidget;
