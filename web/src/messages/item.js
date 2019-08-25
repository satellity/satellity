import style from './item.scss';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import TimeAgo from 'react-timeago';
import Avatar from '../users/avatar.js';

class Item extends Component {
  constructor(props) {
    super(props);
    this.state = props.message;
  }

  render() {
    let state = this.state;
    return (
      <li className={style.message}>
        <div className={style.profile}>
          <Avatar user={state.user} />
          <div>
              {state.user.nickname}
            <div className={style.time}>
              <TimeAgo date={state.created_at} />
            </div>
          </div>
        </div>
        {state.body}
      </li>
    )
  }
}

export default Item;
