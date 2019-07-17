import style from './index.scss';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import API from '../api/index.js';
import New from './new.js';
import Avatar from '../users/avatar.js';

class Index extends Component {
  constructor(props) {
    super(props);
    this.api = new API();

    let id = this.props.match.params.id;
    this.state = {
      group_id: id,
      name: '',
      messages: [],
      loading: true
    };
  }

  componentDidMount() {
    this.api.group.show(this.state.group_id).then((data) => {
      this.setState({name: data.name}, () => {
        this.api.message.index(this.state.group_id).then((data) => {
          this.setState({loading: false, messages: data});
        });
      });
    });
  }

  render() {
    const state = this.state;

    const messages = state.messages.map((message) => {
      return (
        <li key={message.message_id}>
          <div className={style.profile}>
            <Avatar user={message.user} />
            {message.user.nickname}
          </div>
          {message.body}
        </li>
      )
    });

    return (
      <div className='container'>
        <main className='column main'>
          <New groupId={state.group_id} />
          <ul className={style.messages}>
            {messages}
          </ul>
        </main>
        <aside className='column aside'>
          <Link to={`/groups/${state.group_id}`}>
            {state.name} >>
          </Link>
        </aside>
      </div>
    )
  }
}

export default Index;
