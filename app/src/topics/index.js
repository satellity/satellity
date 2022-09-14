import React, {useState, useEffect} from 'react';
import {Link, useParams, useSearchParams} from 'react-router-dom';
import {Helmet} from 'react-helmet';
import Loading from 'components/loading.js';
import API from 'api/index.js';
import Widget from 'home/widget.js';
import {useCategory} from 'services';
import {site} from 'utils';
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


  let clazz = style.node;
  if (categoryId === 'latest') {
    clazz += ` ${style.current}`;
  }
  const latest = (
    <Link to='/' key='latest' className={`${clazz}`}>
      {i18n.t('home.latest')}
    </Link>
  );

  let categories = [latest];
  if (!isLoading) {
    categories = categories.concat(data.map((category) => {
      clazz = style.node;
      if (categoryId === category.name) {
        clazz += ` ${style.current}`;
      }
      return (
        <Link
          to={`/categories/${category.name}`}
          className={`${clazz}`}
          key={category.category_id}>
          {category.alias}
        </Link>
      );
    }));
  }

  return (
    <div className={style.nodes}>
      {categories}
    </div>
  );
};

const Topics = () => {
  const {id} = useParams();
  const [searchParams] = useSearchParams();

  const [loading, setLoading] = useState(true);
  const [categoryId, setCategoryId] = useState(id || 'latest');
  const [pagination] = useState(30);
  const [offset, setOffset] = useState(searchParams.get('offset') || '');
  const [topics, setTopics] = useState([]);

  useEffect(() => {
    setCategoryId(id || 'latest');
    setOffset('');
  }, [id]);

  useEffect(() => {
    setLoading(true);
    const request = api.category.topics(categoryId, offset);

    request.then((resp) => {
      if (resp.error) {
        return;
      }
      console.log(resp);
      const data = resp.data;
      const offset = data.length >= pagination ? data[data.length-1].created_at : '';
      setOffset(offset);
      setTopics(data);
      setLoading(false);
    });
  }, [categoryId, searchParams.get('offset')]);

  if (loading) {
    return (
      <div className={style.loading}>
        <Loading />
      </div>
    );
  }

  const topicsView = topics.map((topic) => {
    return (
      <TopicItem topic={topic} key={topic.topic_id}/>
    );
  });

  return (
    <>
      <ul className={style.topics}> {topicsView} </ul>
      {
        topics.length >= pagination && offset &&
          (
            <div className={style.load}>
              <Link to={`?offset=${offset}`}>{i18n.t('general.next')}</Link>
            </div>
          )
      }
    </>
  );
};

const Index = () => {
  const title = `${i18n.t('site.title')} - ${site.Name}`;
  const description = i18n.t('site.description');

  return (
    <div className='container'>
      <Helmet>
        <title>{title}</title>
        <meta name='description' content={description} />
      </Helmet>
      <main className='column main'>
        <Nodes />
        <Topics />
      </main>
      <aside className='column aside'>
        <Widget />
      </aside>
    </div>
  );
};

export default Index;
