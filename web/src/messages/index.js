import style from './index.scss';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import API from '../api/index.js';
import New from './new.js';
import Item from './item.js';

class Index extends Component {
  constructor(props) {
    super(props);
    this.api = new API();

    let id = this.props.match.params.id;
    this.state = {
      group_id: id,
      name: '',
      messages: {},
      current: {},
      loading: true
    };
  }

  componentDidMount() {
    if (!this.api.user.loggedIn()) {
      return
    }

    this.api.group.show(this.state.group_id).then((data) => {
      this.setState({name: data.name}, () => {
        this.api.message.index(this.state.group_id, '').then((data) => {
          let map = {};
          for (let i=0;i<data.length;i++) {
            let item = data[i];
            if (item.parent_id == item.message_id)  {
              item.children = [];
              map[item.message_id] = item;
            } else {
              if (item[item.parent_id].children) {
                item.children.concat(item)
              }
            }
          }
          this.setState({loading: false, messages: map});
        });
      });
    });
  }

  render() {
    let state = this.state;
    if (!this.api.user.loggedIn()) {
      return (
        <Redirect to={`/groups/${state.group_id}`} />
      )
    }

    let messages = Object.keys(state.messages).map((key) => {
      let message = state.messages[key];
      return (
        <Item message={message}  key={message.message_id} />
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
