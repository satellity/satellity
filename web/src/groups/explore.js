import style from './explore.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class Explore extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {groups: []};
  }

  componentDidMount() {
    this.api.group.index().then((data) => {
      this.setState({groups: data});
    });
  }

  render() {
    const list = this.state.groups.map((group) => {
      let user = group.user;
      return (
        <div key={group.group_id} className={style.item}>
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
        </div>
      )
    });

    return (
      <div className='wrapper container'>
        <div className={style.explore}>
          <h1>
            {i18n.t('group.explore')}
          </h1>
          <div className={style.list}>
            {list}
          </div>
        </div>
      </div>
    )
  }
}

export default Explore;
