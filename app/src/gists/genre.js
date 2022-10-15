import React from 'react';
import {useParams} from 'react-router-dom';
import {useGenres} from 'services';
import Loading from 'components/loading.js';
import Item from './item.js';
import Layout from './layout.js';

import style from './index.module.scss';

const List = () => {
  const {genre} = useParams();
  const {isLoading, data} = useGenres(genre);

  if (isLoading) {
    return (
      <div className={style.loading}>
        <Loading />
      </div>
    );
  }

  const list = data.map((g) => {
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
