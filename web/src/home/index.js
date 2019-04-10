import style from './index.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';
import SiteWidget from './widget.js';
import TopicItem from '../topics/item.js';
import LoadingView from '../loading/loading.js';

class Home extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    this.params = new URLSearchParams(props.location.search);
    this.pagination = 50;
    let categories = [];
    let d = window.localStorage.getItem('categories');
    if (d !== null && d !== undefined && d !== '') {
      categories = JSON.parse(atob(d));
    }
    this.state = {topics: [], categories: categories, category: 'latest', loading: true, offset: ''};

    this.loadTopics = this.loadTopics.bind(this);
  }

  componentDidMount() {
    this.api.category.index().then((data) => {
      this.setState({categories: data});
      let categoryId = 'latest';
      let category = this.params.get('c');
      if (!!category) {
        for (let i = 0; i < this.state.categories.length; i++) {
          let c = this.state.categories[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            categoryId = c.category_id;
            break;
          }
        }
      }
      this.fetchTopics(categoryId, true);
    });
  }

  componentDidUpdate(prevProps) {
    let props = this.props;
    if (props.location.search !== prevProps.location.search) {
      let categoryId = 'latest';
      let params = new URLSearchParams(props.location.search);
      let category = params.get('c');
      if (!!category) {
        // TODO better method?
        for (let i = 0; i < this.state.categories.length; i++) {
          let c = this.state.categories[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            categoryId = c.category_id;
            break;
          }
        }
      }
      this.fetchTopics(categoryId, true);
    }
  }

  loadTopics(e) {
    e.preventDefault();
    this.fetchTopics(this.state.category, false);
  }

  fetchTopics(categoryId, replace) {
    this.setState({loading: replace, offset: ''});
    let request;
    if (categoryId === 'latest') {
      request = this.api.topic.index(this.state.offset);
    } else {
      request = this.api.category.topics(categoryId, this.state.offset);
    }

    request.then((data) => {
      let offset = '';
      if (data.length > this.pagination) {
        offset = data[data.length-1].created_at;
      }
      if (!replace) {
        data = this.state.topics.concat(data);
      }
      this.setState({category: categoryId, loading: false, offset: offset, topics: data});
    });
  }

  render() {
    let state = this.state;

    const loadingView = (
      <div className={style.loading}>
        <LoadingView style='md-ring'/>
      </div>
    )

    const topics = state.topics.map((topic) => {
      return (
        <TopicItem topic={topic} key={topic.topic_id}/>
      )
    });

    const categories = state.categories.map((category) => {
      return (
        <Link
          to={{ pathname: "/", search: `?c=${category.name}` }}
          className={`${style.node} ${state.category === category.category_id ? style.current : ''}`}
          key={category.category_id}>{category.alias}</Link>
      )
    });

    return (
      <div className='container'>
        <main className='section main'>
          <div className={style.nodes}>
            <Link to='/'
              className={`${style.node} ${state.category === 'latest' ? style.current : ''}`}>{i18n.t('home.latest')}</Link>
            {categories}
          </div>
          {state.loading && loadingView}
          {!state.loading && <ul className={style.topics}> {topics} </ul>}
          {state.offset !== '' && <div className={style.load}><a href='javascript:;' onClick={this.loadTopics}>Load More</a></div>}
        </main>
        <aside className='section aside'>
          <SiteWidget />
        </aside>
      </div>
    );
  }
}

export default Home;
