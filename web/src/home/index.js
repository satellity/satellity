import style from './index.scss';
import topicStyle from '../styles/topic_item.scss';
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
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (d !== null && d !== undefined && d !== '') {
      categories = JSON.parse(atob(d));
    }
    this.state = {topics: [], categories: categories, category: 'latest'};
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('home', 'layout');
    this.handleClick = this.handleClick.bind(this);
  }

  componentDidMount() {
    this.api.topic.index((resp) => {
      this.setState({topics: resp.data});
    });
    this.api.category.index((resp) => {
      this.setState({categories: resp.data});
    });
  }

  handleClick(id, e) {
    e.preventDefault();
    this.setState({category: id});
    if (id === 'latest') {
      this.api.topic.index((resp) => {
        this.setState({topics: resp.data});
      });
      return
    }
    this.api.category.topics(id, (resp) => {
      this.setState({topics: resp.data});
    });
  }

  render() {
    return (
      <HomeView state={this.state} color={this.color} handleClick={this.handleClick} />
    );
  }
}

const HomeView = (props) => {
  const topics = props.state.topics.map((topic) => {
    let comment = '';
    if (topic.comments_count > 0) {
      comment = (
        <span className={topicStyle.count} style={{backgroundColor: props.color.colour(topic.topic_id)}}> {topic.comments_count} </span>
      )
    }
    return (
      <li className={topicStyle.topic} key={topic.topic_id}>
        <img src={topic.user.avatar_url} className={topicStyle.avatar} />
        <div className={topicStyle.detail}>
          <h2 className={topicStyle.title}>
            <Link to={`/topics/${topic.topic_id}`}>{topic.title}</Link>
          </h2>
          <span>{topic.category.name}</span> • {topic.user.nickname} • <TimeAgo date={topic.created_at} />
        </div>
        <div className={topicStyle.comment}>
          {comment}
        </div>
      </li>
    )
  });

  const categories = props.state.categories.map((category) => {
    return (
      <Link
        to="/"
        className={`${style.node} ${props.state.category === category.category_id ? style.current : ''}`}
        onClick={(e) => props.handleClick(category.category_id, e)}
        key={category.category_id}>{category.name}</Link>
    )
  });

  return (
    <div className='container'>
      <main className='section main'>
        <div className={style.nodes}>
          <Link to='/'
            className={`${style.node} ${props.state.category === 'latest' ? style.current : ''}`}
            onClick={(e) => props.handleClick('latest', e)}>Latest</Link>
          {categories}
        </div>
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
