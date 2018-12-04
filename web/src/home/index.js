import './index.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Typed from 'typed.js';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';

class Home extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {topics: []};
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('home', 'layout');
  }

  componentDidMount() {
    this.api.topic.index((resp) => {
      this.setState({topics: resp.data});
    });
  }

  render() {
    return (
      <HomeView state={this.state} />
    );
  }
}

const HomeView = (props) => {
  const topics = props.state.topics.map((topic) => {
    return (
      <li className='topic item' key={topic.topic_id}>
        <img src={topic.user.avatar_url} className='avatar' />
        <div className='detail'>
          <h2>
            <Link to={`/topics/${topic.topic_id}`}>{topic.title}</Link>
            <span className='comment'>
              <FontAwesomeIcon icon={['far', 'comment']} />
              {topic.comments_count}
            </span>
          </h2>
          <span className='category'>{topic.category.name}</span> • {topic.user.nickname} • <TimeAgo date={topic.created_at} />
        </div>
      </li>
    )
  });

  return (
    <div className='app container'>
      <main className='section main'>
        <ul>
          {topics}
        </ul>
      </main>
      <aside className='section aside'>
        <div className='section site info'>
          SUNTIN is an open source forum, get codes in <Link to='https://github.com/godiscourse/godiscourse'>Github</Link>.
        </div>
      </aside>
    </div>
  );
}

export default Home;
