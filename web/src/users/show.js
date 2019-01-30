import style from './show.scss';
import moment from 'moment';
import React, { Component } from 'react';
import API from '../api/index.js';

class UserShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {user_id: props.match.params.id, nickname: '', biography: '', avatar_url: '', created_at: ''}
  }

  componentDidMount() {
    this.api.user.show(this.state.user_id, (resp) => {
      let user = resp.data;
      user.created_at = moment(user.created_at).format('l');
      this.setState(user);
    });
  }

  render() {
    return (
      <View state={this.state} />
    )
  }
}

const View = (props) => {
  return (
    <div className='container'>
      <aside className='section aside'>
        <div className={style.profile}>
          <img src={props.state.avatar_url} className={style.avatar} />
          <div className={style.name}>
            {props.state.nickname}
          </div>
          <div className={style.created}>
            Joined {props.state.created_at}
          </div>
        </div>
      </aside>
      <main className='section main'>
      </main>
    </div>
  )
};

export default UserShow;
