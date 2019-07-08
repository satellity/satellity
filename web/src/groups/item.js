import style from './item.scss';
import React, { Component } from 'react';
import {Link} from 'react-router-dom';

class Item extends Component {
  constructor(props) {
    super(props);
    this.state = props.group;
  }

  render() {
    let group = this.state;
    let user = group.user;
    return (
      <div className={style.group}>
        <div className={style.profile}>
          <img src={user.avatar_url} alt={user.nickname} className={style.avatar} />
          <div className={style.nickname}>
            {user.nickname}
          </div>
          <div>
            {group.users_count}
          </div>
        </div>
        <h2 className={style.name}>
          <Link to={`/groups/${group.group_id}`}>{group.name}</Link>
        </h2>
        <div>
          {group.description.slice(0, 120)}
        </div>
      </div>
    )
  }
}

export default Item;
