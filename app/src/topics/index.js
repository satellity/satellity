import React, {useState, useEffect} from 'react';
import {Link} from 'react-router-dom';
import {Helmet} from 'react-helmet';
import Config from 'components/config.js';
import Loading from 'components/loading.js';
import API from 'api/index.js';
import Widget from 'home/widget.js';
import TopicItem from './item.js';

import style from './index.module.scss';

const api = new API();

const Index = () => {
  const [i18n] = useState(window.i18n);
  const [loading, setLoading] = useState(true);
  const [categoryId] = useState('latest');
  const [pagination] = useState(30);
  const [offset, setOffset] = useState('');
  const [topics, setTopics] = useState([]);
  const [categories] = useState([]);

  useEffect(() => {
    const request = categoryId === 'latest' ?
      api.topic.index(offset) :
      api.category.topics(category_id, offset);

    request.then((resp) => {
      if (resp.error) {
        return;
      }
      const data = resp.data;
      const offset = data.length >= pagination ? data[data.length-1].created_at : '';
      setOffset(offset);
      setTopics(data);
      setLoading(false);
    });
  }, []);

  const loadingView = (
    <div className={style.loading}>
      <Loading />
    </div>
  );

  const topicsView = topics.map((topic) => {
    return (
      <TopicItem topic={topic} key={topic.topic_id}/>
    );
  });

  const categoriesView = categories.map((category) => {
    return (
      <Link to={`/categories/${category.name}`} className={`${style.node} ${state.category_id === category.category_id ? style.current : ''}`}
        key={category.category_id}>
        {category.alias}
      </Link>
    );
  });

  const title = `${i18n.t('site.title')} - ${Config.Name}`;
  const description = i18n.t('site.description');
  const canonical = <link rel="canonical" href={`${Config.Host}`} />;

  return (
    <div className='container'>
      {
        !loading &&
          <Helmet>
            <title>{title}</title>
            <meta name='description' content={description} />
            {canonical}
          </Helmet>
      }
      <main className='column main'>
        <div className={style.nodes}>
          <Link to='/'
            className={`${categoryId === 'latest' ? style.current : ''}`}>{i18n.t('home.latest')}</Link>
          {categoriesView}
        </div>

        {loading && loadingView}

        {!loading && <ul className={style.topics}> {topicsView} </ul>}
        {
          topics.length >= pagination && offset &&
            (
              <div className={style.load}>
                <Link to={`?offset=${offset}`}>{i18n.t('general.next')}</Link>
              </div>
            )
        }
      </main>
      <aside className='column aside'>
        <Widget />
      </aside>
    </div>
  );
};

export default Index;
