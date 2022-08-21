import style from './index.module.scss';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import {Helmet} from 'react-helmet';
import Config from '../components/config.js';
import Base64 from '../components/base64.js';
import Loading from '../components/loading.js';
import API from '../api/index.js';
import TopicItem from './item.js';
import Widget from '../home/widget.js';

class Index extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
    this.base64 = new Base64();
    // TODO decode should categories in api;
    let categories = [];
    const d = window.localStorage.getItem('categories');
    if (d !== null && d !== undefined && d !== '') {
      categories = JSON.parse(this.base64.decode(d));
    }
    this.state = {
      id: 'latest',
      user: {},
      topics: [],
      categories: categories,
      category: {},
      category_id: 'latest',
      loading: true,
      pagination: 30,
      offset: 'TODO',
    };
  }

  componentDidMount() {
    const user = this.api.user.local();
    this.api.category.index().then((resp) => {
      if (resp.error) {
        return;
      }
      let categoryId = 'latest';
      let current = {};
      const category = this.state.id;
      if (category !== 'latest') {
        for (let i = 0; i < resp.data.length; i++) {
          const c = resp.data[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            categoryId = c.category_id;
            current = c;
            break;
          }
        }
      }
      this.setState({user: user, category: current, category_id: categoryId, categories: resp.data}, () => {
        this.fetchTopics(categoryId);
      });
    });
  }

  componentDidUpdate(prevProps) {
    if (false) {
      let categoryId = 'latest';
      let current = {};
      const category = 'latest';
      if (category !== 'latest') {
        for (let i = 0; i < this.state.categories.length; i++) {
          const c = this.state.categories[i];
          if (c.name.toLocaleLowerCase() === category.toLocaleLowerCase()) {
            categoryId = c.category_id;
            current = c;
            break;
          }
        }
      }
      this.setState({category: current, category_id: categoryId, offset: 'TODO'}, () => {
        this.fetchTopics(category_id);
      });
    }
  }

  fetchTopics(categoryId) {
    this.setState({loading: true}, () => {
      const request = categoryId === 'latest' ?
        this.api.topic.index(this.state.offset) :
        this.api.category.topics(category_id, this.state.offset);

      request.then((resp) => {
        if (resp.error) {
          return;
        }
        const data = resp.data;
        const offset = data.length >= this.state.pagination ? data[data.length-1].updated_at : '';
        this.setState({category_id: categoryId, loading: false, offset: offset, topics: data});
      });
    });
  }

  render() {
    const i18n = window.i18n;
    const state = this.state;

    const loadingView = (
      <div className={style.loading}>
        <Loading />
      </div>
    );

    const topics = state.topics.map((topic) => {
      return (
        <TopicItem user={state.user} topic={topic} key={topic.topic_id}/>
      );
    });

    const categories = state.categories.map((category) => {
      return (
        <Link to={`/categories/${category.name}`} className={`${style.node} ${state.category_id === category.category_id ? style.current : ''}`}
          key={category.category_id}>
          {category.alias}
        </Link>
      );
    });

    let title; let description;
    if (!!state.category.name) {
      title = `${state.category.alias} - ${Config.Name}`;
      description = state.category.description;
    } else {
      title = `${i18n.t('site.title')} - ${Config.Name}`;
      description = i18n.t('site.description');
    }

    let canonical = <link rel="canonical" href={`${Config.Host}`} />;
    if (state.category_id !== 'latest') {
      canonical = <link rel="canonical" href={`${Config.Host}/categories/${state.category.name}`} />;
    }

    return (
      <div className='container'>
        {
          !state.loading &&
          <Helmet>
            <title>{title}</title>
            <meta name='description' content={description} />
            {canonical}
          </Helmet>
        }
        <main className='column main'>
          <div className={style.nodes}>
            <Link to='/'
              className={`${style.node} ${state.category_id === 'latest' ? style.current : ''}`}>{i18n.t('home.latest')}</Link>
            {categories}
          </div>
          {state.loading && loadingView}
          {!state.loading && <ul className={style.topics}> {topics} </ul>}
          {
            state.topics.length >= state.pagination && state.offset &&
              (
                <div className={style.load}>
                  <Link to={`${props.match.url}?offset=${state.offset}`}>{i18n.t('general.next')}</Link>
                </div>
              )
          }
        </main>
        <aside className='column aside'>
          <Widget />
        </aside>
      </div>
    );
  }
}

export default Index;
