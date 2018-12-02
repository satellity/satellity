import './index.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';

class TopicShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {
      topic_id: '', title: '', body: '', comments_count: 0, created_at: '', user_id: '', is_author: false,
      user: {user_id: '', nickname: ''},
      category: {category_id: '', name: ''}
    };
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('topic', 'layout');
  }

  componentDidMount() {
    const user = this.api.user.me();
    this.api.topic.show(this.props.match.params.id, (resp) => {
      resp.data.is_author = resp.data.user.user_id === user.user_id;
      this.setState(resp.data);
    });
  }

  render() {
    return (
      <View state={this.state} />
    )
  }
}

const View = ({state}) => {
  var editAction = '';
  if (state.is_author) {
    editAction = (
      <Link to={`/topics/${state.topic_id}/edit`}>
        <FontAwesomeIcon icon={['far', 'edit']} />
      </Link>
    )
  }
  return (
    <div>
      <header className='topic header'>
        <h1>
          {state.title}
          {editAction}
          <img src={state.user.avatar_url} className='avatar' />
        </h1>
        <div className='info'>
          {state.category.name} • {state.user.nickname} • <TimeAgo date={state.created_at} />
        </div>
      </header>
      <div className='topic body'>
        {state.body}
      </div>
    </div>
  )
}

export default TopicShow;
