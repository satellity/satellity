import style from './index.scss';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import API from '../api/index.js';
import New from './new.js';
import Item from './item.js';

class Index extends Component {
  constructor(props) {
    super(props);
    let id = this.props.match.params.id;
    this.state = {
      group_id: id,
      name: '',
      messages: [],
      current: {},
      loading: true
    };

    this.api = new API();
    this.handleCommentClick = this.handleCommentClick.bind(this);
  }

  componentDidMount() {
    if (!this.api.user.loggedIn()) {
      return
    }

    this.api.group.show(this.state.group_id).then((data) => {
      this.setState({name: data.name}, () => {
        this.api.message.index(this.state.group_id, '').then((data) => {
          let array = [];
          let mid = {children: []};
          for (let i=0;i<data.length;i++) {
            let item = data[i];
            if (item.parent_id == item.message_id)  {
              item.children = mid.children;
              array.push(item);
              mid.children = [];
            } else {
              mid.children.push(item);
            }
          }
          this.setState({loading: false, messages: array});
        });
      });
    });
  }

  handleCommentClick(id) {
    this.setState({current: {message_id: id}});
  }

  render() {
    let state = this.state;
    if (!this.api.user.loggedIn()) {
      return (
        <Redirect to={`/groups/${state.group_id}`} />
      )
    }

    let messages = state.messages.map((message) => {
      return (
        <Item message={message} current={state.current} key={message.message_id} handleComment={this.handleCommentClick} />
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
