import Style from './index.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
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
      <div>
        <h1>
          {i18n.t('site.slogan')}
        </h1>
        <div className={Style.features}>
          <div className={Style.section}>
            <FontAwesomeIcon icon={['fa', 'chalkboard']} />
            <div className={Style.desc}>
              {i18n.t('home.forum')}
            </div>
          </div>
          <div className={Style.section}>
            <FontAwesomeIcon icon={['fa', 'users-cog']} />
            <div className={Style.desc}>
              {i18n.t('home.group')}
            </div>
          </div>
        </div>
        <div>
        </div>
      </div>
    )
  }
}

export default Index;
