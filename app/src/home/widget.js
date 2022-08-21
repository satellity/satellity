import React from 'react';
import API from 'api/index.js';
import Button from 'components/button.js';

import style from './widget.module.scss';

const SiteWidget = () => {
  const api = new API();
  const i18n = window.i18n;

  return (
    <div className={style.widget}>
      {
        api.user.loggedIn() && (
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
      <div className={style.copyright}>
        © 2022 - Now MIT license
      </div>
    </div>
  );
};

export default SiteWidget;
