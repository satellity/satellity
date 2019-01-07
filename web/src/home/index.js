import './index.scss';
import style from './index.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Typed from 'typed.js';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import TimeAgo from 'react-timeago';
import API from '../api/index.js';
import ColorUtils from '../components/color.js';
import SiteWidget from '../components/site-widget.js';

class Home extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.color = new ColorUtils();
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
      <HomeView state={this.state} color={this.color} />
    );
  }
}

const HomeView = (props) => {
  const topics = props.state.topics.map((topic) => {
    let comment = '';
    if (topic.comments_count > 0) {
      comment = (
        <div className={style.comment}>
          <span className={style.count} style={{backgroundColor: props.color.colour(topic.topic_id)}}> {topic.comments_count} </span>
        </div>
      )
    }
    return (
      <li className={style.topic} key={topic.topic_id}>
        <img src={topic.user.avatar_url} className={style.avatar} />
        <div className={style.detail}>
          <h2 className={style.title}>
            <Link to={`/topics/${topic.topic_id}`}>{topic.title}</Link>
          </h2>
          <span>{topic.category.name}</span> • {topic.user.nickname} • <TimeAgo date={topic.created_at} />
        </div>
        {comment}
      </li>
    )
  });

  return (
    <div className='container'>
      <main className='section main'>
        <ul className={style.topics}>
          {topics}
        </ul>
      </main>
      <aside className='section aside'>
        <SiteWidget />
      </aside>
    </div>
  );
}

export default Home;
