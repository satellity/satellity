import style from './item.scss';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import TimeAgo from 'react-timeago';
import Avatar from '../users/avatar.js';

class Item extends Component {
  constructor(props) {
    super(props);
    this.state = {
      message: props.message,
      current: props.current,
    }
  }

  render() {
    let state = this.state;
    return (
      <li className={style.message}>
        <div className={style.profile}>
          <Avatar user={state.message.user} />
          <div>
              {state.message.user.nickname}
            <div className={style.time}>
              <TimeAgo date={state.message.created_at} />
            </div>
          </div>
        </div>
        {state.message.body}
      </li>
    )
  }
}

export default Item;
