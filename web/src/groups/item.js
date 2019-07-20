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
        <div className={style.head}>
          <div className={style.title}>
            <h2 className={style.name}>
              <Link to={`/groups/${group.group_id}`}>{group.name}</Link>
            </h2>
            <div className={style.nickname}>
              OWNER {user.nickname}
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
