import style from './profile.module.scss';
import React, {Component} from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class Profile extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = this.api.user.local();
  }

  componentDidMount() {
    this.api.user.remote().then((resp) => {
      if (resp.error) {
        return
      }
      this.setState(resp.data);
    });
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;
    return (
      <div className={style.profile}>
        <div className={style.user}>
          <img src={state.avatar_url} className={style.avatar} alt={state.nickname} />
          <div className={style.nickname}>
              {state.nickname}
          </div>
        </div>
        <div className={style.new}>
          <Link to='/topics/new'>{i18n.t('topic.new')}</Link>
        </div>
      </div>
    )
  }
}

export default Profile;
