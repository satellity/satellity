import React, {useState, useEffect} from 'react';
import {useSearchParams} from 'react-router-dom';
import API from 'api/index.js';
import Loading from 'components/loading.js';
import Item from './item.js';
import Layout from './layout.js';

import style from './index.module.scss';

const api = new API();

const List = () => {
  const [searchParams] = useSearchParams();

  const [loading, setLoading] = useState(true);
  const [gists, setGists] = useState([]);
  const [offset, setOffset] = useState(searchParams.get('offset') || '');
  const [pagination] = useState(128);

  useEffect(() => {
    setLoading(true);
    const request = api.gist.index(offset);

    request.then((resp) => {
      if (resp.error) {
        return;
      }
      const data = resp.data;
      const offset = data.length >= pagination ? data[data.length-1].created_at : '';
      setOffset(offset);
      setGists(data);
      setLoading(false);
    });
  }, [searchParams.get('offset')]);

  if (loading) {
    return (
      <div className={style.loading}>
        <Loading />
      </div>
    );
  }

  const list = gists.map((g) => {
    return (
      <Item key={g.gist_id} gist={g} />
    );
  });

  return (
    <>
      {list}
    </>
  );
};

const Index = () => {
  return (
    <Layout>
      <List />
    </Layout>
  );
};

export default Index;
