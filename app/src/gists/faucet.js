import React from 'react';
import Loading from 'components/loading.js';
import {useFaucet} from 'services';
import Layout from './layout.js';

import style from './faucet.module.scss';

const List = () => {
  const {isLoading, data} = useFaucet();

  if (isLoading) {
    return (
      <div className={style.loading}>
        <Loading />
      </div>
    );
  }

  const views = data.filter((d) => {
    if (d.chainId === 3) {
      return false;
    }
    if (d.chainId === 4) {
      return false;
    }
    return d.faucets.length > 0;
  }).map((d) => {
    const fs = d.faucets.map((f, i) => {
      const link = new URL(f);
      return (
        <div key={f}>
          <a href={f}> {i}. {link.host} 🔗 </a>
        </div>
      );
    });

    return (
      <div key={d.chainId} className={style.item}>
        <div className={style.name}>
          {d.name} Faucet
        </div>
        <div>
          {fs}
        </div>
      </div>
    );
  });

  return (
    <>
      {views}
    </>
  );
};

const Faucet = () => {
  return (
    <Layout>
      <List />
    </Layout>
  );
};

export default Faucet;
