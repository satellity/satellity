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
import LoadingView from '../loading/loading.js';

class Home extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    this.color = new ColorUtils();
    this.params = new URLSearchParams(props.location.search);
    this.pagination = 50;
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (d !== null && d !== undefined && d !== '') {
      categories = JSON.parse(atob(d));
    }
    this.state = {topics: [], categories: categories, category: 'latest', loading: true, offset: ''};
    this.handleClick = this.handleClick.bind(this);
    this.load = this.load.bind(this);
  }

  componentDidMount() {
    this.api.category.index().then((data) => {
      this.setState({categories: data});
    })
    const category = this.params.get("c");
    let request, id = 'latest';
    if (!!category) {
      for (let i=0; i< this.state.categories.length; i++) {
        let c = this.state.categories[i];
        if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
          id = c.category_id, request = this.api.category.topics(id);
          break;
        }
      }
    }
    if (id === 'latest') {
      request = this.api.topic.index();
    }
    request.then((data) => {
      let offset = '';
      if (data.length === this.pagination) {
        offset = data[data.length-1].created_at;
      }
      this.setState({category: id, topics: data, loading: false, offset: offset});
    });
  }

  handleClick(id, e) {
    e.preventDefault();
    this.setState({loading: true, offset: ''});
    let request;
    if (id === 'latest') {
      request = this.api.topic.index();
    } else {
      request = this.api.category.topics(id)
    }
    request.then((data) => {
      let offset = '';
      if (data.length == this.pagination) {
        offset = data[data.length-1].created_at;
      }
      this.setState({category: id, topics: data, loading: false, offset: offset});
    });
  }

  load(e) {
    e.preventDefault();
    let id = this.state.category;
    let request;
    if (id === 'latest') {
      request = this.api.topic.index(this.state.offset);
    } else {
      request = this.api.category.topics(id, this.state.offset);
    }
    request.then((data) => {
      let offset = '';
      if (data.length === this.pagination) {
        offset = data[data.length-1].created_at;
      }
      data = this.state.topics.concat(data);
      this.setState({category: id, topics: data, loading: false, offset: offset});
    });
  }

  render() {
    return (
      <HomeView state={this.state} color={this.color} handleClick={this.handleClick} load={this.load} />
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
        key={category.category_id}>{category.alias}</Link>
    )
  });

  const loadingView = (
    <div className={style.loading}>
      <LoadingView style='md-ring'/>
    </div>
  )

  return (
    <div className='container'>
      <main className='section main'>
        <div className={style.nodes}>
          <Link to='/'
            className={`${style.node} ${props.state.category === 'latest' ? style.current : ''}`}
            onClick={(e) => props.handleClick('latest', e)}>Latest</Link>
          {categories}
        </div>
        {props.state.loading && loadingView}
        <ul className={style.topics}>
          {topics}
        </ul>
        {props.state.offset !== '' && <div className={style.load}><a href='javascript:;' onClick={(e) => props.load(e)}>Load More</a></div>}
      </main>
      <aside className='section aside'>
        <SiteWidget />
      </aside>
    </div>
  );
}

export default Home;
