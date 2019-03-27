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
    this.loadMore = this.loadMore.bind(this);
  }

  componentDidMount() {
    this.api.category.index().then((data) => {
      this.setState({categories: data});
      let categoryId = 'latest';
      const category = this.params.get("c");
      if (!!category) {
        for (let i=0; i< this.state.categories.length; i++) {
          let c = this.state.categories[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            categoryId = c.category_id;
            break;
          }
        }
      }
      this.loadTopics(categoryId, true);
    });
  }

  componentDidUpdate(prevProps) {
    let props = this.props;
    if (props.location.search !== prevProps.location.search) {
      let categoryId = 'latest';
      let params = new URLSearchParams(props.location.search);
      let category = params.get("c");
      if (!!category) {
        for (let i=0; i< this.state.categories.length; i++) {
          let c = this.state.categories[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            categoryId = c.category_id;
            break;
          }
        }
      }
      console.log(this.state.categories, categoryId);
      this.loadTopics(categoryId, true);
    }
  }

  loadMore(e) {
    e.preventDefault();
    this.loadTopics(this.state.category, false);
  }

  loadTopics(categoryId, reload) {
    this.setState({loading: true, offset: ''});
    let request;
    if (categoryId === 'latest') {
      request = this.api.topic.index(this.state.offset);
    } else {
      request = this.api.category.topics(categoryId, this.state.offset);
    }
    request.then((data) => {
      let offset = '';
      if (data.length === this.pagination) {
        offset = data[data.length-1].created_at;
      }
      if (!reload) {
        data = this.state.topics.concat(data);
      }
      this.setState({category: categoryId, topics: data, loading: false, offset: offset});
    });
  }

  render() {
    return (
      <HomeView state={this.state} color={this.color} loadMore={this.loadMore} />
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
        to={{ pathname: "/", search: `?c=${category.name}` }}
        className={`${style.node} ${props.state.category === category.category_id ? style.current : ''}`}
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
            className={`${style.node} ${props.state.category === 'latest' ? style.current : ''}`}>Latest</Link>
          {categories}
        </div>
        {props.state.loading && loadingView}
        <ul className={style.topics}>
          {topics}
        </ul>
        {props.state.offset !== '' && <div className={style.load}><a href='javascript:;' onClick={(e) => props.loadMore(e)}>Load More</a></div>}
      </main>
      <aside className='section aside'>
        <SiteWidget />
      </aside>
    </div>
  );
}

export default Home;
