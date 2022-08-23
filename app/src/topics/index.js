import React, {useState, useEffect} from 'react';
import {Link} from 'react-router-dom';
import {Helmet} from 'react-helmet';
import {useParams, useSearchParams} from 'react-router-dom';
import Config from 'components/config.js';
import Loading from 'components/loading.js';
import API from 'api/index.js';
import Widget from 'home/widget.js';
import {useCategory} from 'services';
import TopicItem from './item.js';

import style from './index.module.scss';

const api = new API();

const Nodes = () => {
  const {id} = useParams();
  const [categoryId, setCategoryId] = useState(id || 'latest');

  const {isLoading, data} = useCategory();

  useEffect(() => {
    setCategoryId(id || 'latest');
  }, [id]);

  if (isLoading) {
    return;
  }

  const categoriesView = data.map((category) => {
    return (
      <Link to={`/categories/${category.name}`} className={`${style.node} ${categoryId === category.name ? style.current : ''}`}
        key={category.category_id}>
        {category.alias}
      </Link>
    );
  });

  return (
    <div className={style.nodes}>
      <Link to='/' className={`${style.node} ${categoryId === 'latest' ? style.current : ''}`}>
        {i18n.t('home.latest')}
      </Link>
      {categoriesView}
    </div>
  );
};

const Index = () => {
  const {id} = useParams();
  const [searchParams] = useSearchParams();

  const [i18n] = useState(window.i18n);
  const [loading, setLoading] = useState(true);
  const [categoryId, setCategoryId] = useState(id || 'latest');
  const [pagination] = useState(30);
  const [offset, setOffset] = useState(searchParams.get('offset') || '');
  const [topics, setTopics] = useState([]);

  useEffect(() => {
    setCategoryId(id || 'latest');
  }, [id]);

  useEffect(() => {
    setLoading(true);
    const request = categoryId === 'latest' ?
      api.topic.index(offset) :
      api.category.topics(categoryId, offset);

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
  }, [categoryId, searchParams.get('offset')]);

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

  const title = `${i18n.t('site.title')} - ${Config.Name}`;
  const description = i18n.t('site.description');

  return (
    <div className='container'>
      {
        !loading &&
          <Helmet>
            <title>{title}</title>
            <meta name='description' content={description} />
          </Helmet>
      }
      <main className='column main'>
        <Nodes />

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
