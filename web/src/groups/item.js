import style from './item.scss';
import React, { Component } from 'react';
import {Link} from 'react-router-dom';

class Item extends Component {
  constructor(props) {
    super(props);

    this.state = {
      group: props.group,
      to: props.to
    }
  }

  render() {
    let group = this.state.group;
    let user = group.user;
    let membersView = group.users_count>1 ? <span className={style.count}>+{group.users_count-1 }</span> : '';
    let link = `/groups/${group.group_id}`;
    if (this.state.to == 'messages') {
      link = `/groups/${group.group_id}/messages`
    }

    return (
      <div className={style.group}>
        <div className={style.head}>
          <div className={style.title}>
            <h2 className={style.name}>
              <Link to={link}>{group.name}</Link>
            </h2>
            <div className={style.nickname}>
              By {user.nickname}
              {membersView}
            </div>
          </div>
          <img src={user.avatar_url} alt={user.nickname} className={style.avatar} />
        </div>
        <div className={style.description}>
          {group.description.slice(0, 120)}
        </div>
      </div>
    )
  }
}

export default Item;
