import style from './index.module.scss';
import React, { Component } from 'react';
import API from '../../api/index.js';

class AdminUser extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {users: []};
  }

  componentDidMount() {
    this.api.user.admin.index().then((resp) => {
      if (resp.error) {
        return;
      }
      this.setState({users: resp.data});
    });
  }

  render() {
    return (
      <UserIndex users={this.state.users} />
    )
  }
}

const UserIndex = (props) => {
  const listUsers = props.users.map((user) => {
    return (
      <li key={user.user_id} className={style.user}>
        <img src={user.avatar_url} alt={user.nickname} className={style.avatar} />
        <div className={style.detail}>
          <div>
            {user.nickname}
          </div>
          <div className={style.time}>
            {user.user_id} | {user.created_at}
          </div>
        </div>
      </li>
    )
  });

  return (
    <div>
      <h1 className='welcome'>
        The list of registered users.
      </h1>
      <div className='panel'>
        <ul>
          {listUsers}
        </ul>
      </div>
    </div>
  );
}

export default AdminUser;
