import './index.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';

class TopicShow extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {
      topic_id: '', title: '', body: '', comments_count: 0, created_at: '',
      user: {user_id: '', nickname: ''},
      category: {category_id: '', name: ''}
    };
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('topic', 'layout');
  }

  componentDidMount() {
    this.api.topic.show(this.props.match.params.id, (resp) => {
      console.info(resp.data);
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
  return (
    <div>
      <header className='topic header'>
        <h1>
          {state.title}
          <img src={state.user.avatar_url} className='avatar' />
        </h1>
        <div>
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
