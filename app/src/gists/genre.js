import React, {useState, useEffect} from 'react';
import {useParams, useSearchParams} from 'react-router-dom';
import API from 'api/index.js';
import Loading from 'components/loading.js';
import Widget from 'components/widget.js';
import Item from './item.js';
import Chart from './chart.js';

import style from './index.module.scss';

const api = new API();

const List = () => {
  const {genre} = useParams();
  const [searchParams] = useSearchParams();
  console.log(searchParams);

  const [loading, setLoading] = useState(true);
  const [gists, setGists] = useState([]);
  const [offset, setOffset] = useState(searchParams.get('offset') || '');
  const [pagination] = useState(128);

  useEffect(() => {
    setLoading(true);
    const request = api.gist.genres(genre, offset);

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
  }, [genre, offset]);

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
    <div className='container'>
      <main className='column main'>
        <List />
      </main>
      <aside className='column aside'>
        <Widget>
          <Chart />
        </Widget>
      </aside>
    </div>
  );
};

export default Index;
