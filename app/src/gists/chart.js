import React, {useState, useEffect} from 'react';
import API from 'api/index.js';

import style from './index.module.scss';

const api = new API();

const Chart = () => {
  const set = {
    bitcoin: 'Bitcoin',
    ethereum: 'Ethereum',
    binancecoin: 'BNB',
    dogecoin: 'Dogecoin',
  };

  const [tokens, setTokens] = useState([]);

  useEffect(() => {
    api.cgc.coins().then((resp) => {
      const tokens = [];
      const data = resp.data;
      Object.keys(resp.data).forEach((key) => {
        data[key].name = set[key];
        tokens.push(data[key]);
      });
      setTokens(tokens);
    }).catch((err) => {
      console.log('catch', err);
    });
  }, []);

  if (tokens.length < 1) {
    return <></>;
  }

  const views = tokens.sort((i, j) => {
    if (i.usd > j.usd) {
      return -1;
    }
    if (j.usd > i.usd) {
      return 1;
    }
    return 0;
  }).map((token) => {
    const change = `${Number.parseFloat(token.usd_24h_change).toFixed(2)}%`;
    return (
      <div key={token.name} className={style.item}>
        <div>
          <div> {token.name} </div>
          <div className={`${style.second} ${style.unit}`}> USD </div>
        </div>
        <div className={style.price}>
          <div className={style.amount}> {token.usd} </div>
          <div className={`${style.second} ${token.usd_24h_change < 0 ? style.descent : style.ascent}`}>${change}</div>
        </div>
      </div>
    );
  });

  return (
    <div className={style.box}>
      {views}
    </div>
  );
};

export default Chart;
