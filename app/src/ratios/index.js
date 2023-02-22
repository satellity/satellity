import React from 'react';
import BigNumber from 'bignumber.js';
import Loading from 'components/loading.js';
import {useRatios} from 'services';
import Widget from 'components/widget.js';
import style from './index.module.scss';

const Index = () => {
  const {isLoading, data} = useRatios();

  if (isLoading) {
    return (
      <div className={style.loading}>
        <Loading />
      </div>
    );
  }

  const list = data.sort((a, b) => {
    const r = new BigNumber(a.global_ratio);
    return r.comparedTo(b.global_ratio) * -1;
  }).map((r) => {
    const ratios = r.ratios.map((rr) => {
      return (
        <div key={rr.category}>
          {rr.long_short_ratio} / {rr.long_account} / {rr.short_account}
        </div>
      );
    });
    return (
      <div key={r.symbol} className={style.item}>
        <div className={style.info}>
          <img src={r.image} alt={r.name} className={style.icon} />
          {r.name} / {r.symbol}
        </div>
        <div>
          {r.current_price} / {r.high_24h} / {r.low_24h}
        </div>
        <div>
          {r.market_cap_rank} / {r.market_cap}
        </div>
        {ratios}
      </div>
    );
  });

  return (
    <div className='container'>
      <main className='column main'>
        {list}
      </main>
      <aside className='column aside'>
        <Widget />
      </aside>
    </div>
  );
};

export default Index;
