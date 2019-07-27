import style from './list.scss';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import API from '../api/index.js';
import GroupItem from './item.js';

class List extends Component {
  constructor(props) {
    super(props);
    this.api = new API();

    this.state = {
      isLoggedIn: this.api.user.loggedIn(),
      groups: []
    }
  }

  componentDidMount() {
    if (this.state.isLoggedIn) {
      this.api.me.groups(90).then((data) => {
        this.setState({groups: data});
      });
    }
  }

  render() {
    const state = this.state;

    if (!state.isLoggedIn) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }
    let list = state.groups.map((group) => {
      return (
        <div key={group.group_id} className={style.item}>
          <GroupItem group={group} to='messages' />
        </div>
      )
    });
    return (
      <div className='wrapper container'>
        <div className={style.panel}>
          <h1>
            {i18n.t('group.dashboard')}
          </h1>
          <div className={style.groups}>
            {list}
          </div>
        </div>
      </div>
    )
  }
}

export default List;
