import React, {useEffect, useState} from 'react';
import BigNumber from 'bignumber.js';
import Loading from 'components/loading.js';
import {useRatios} from 'services';
import Widget from 'components/widget.js';
import style from './index.module.scss';

const Index = () => {
  const [ratiosData, setRatiosData] = useState([]);
  const [arrow, setArrow] = useState('↕');
  const {isLoading, data} = useRatios();

  useEffect(() => {
    if (!isLoading) {
      setRatiosData(data.sort((a, b) => {
        const r = new BigNumber(a.market_cap_rank);
        return r.comparedTo(b.market_cap_rank);
      }));
    }
  }, [isLoading]);

  if (isLoading) {
    return (
      <div className={style.loading}>
        <Loading />
      </div>
    );
  }

  const handleSort = (e) => {
    if (arrow === '↓') {
      setArrow('↑');
      setRatiosData(ratiosData.sort((a, b) => {
        const r = new BigNumber(a.global_ratio);
        return r.comparedTo(b.global_ratio);
      }));
      return;
    }
    if (arrow === '↑') {
      setArrow('↕');
      setRatiosData(ratiosData.sort((a, b) => {
        const r = new BigNumber(a.market_cap_rank);
        return r.comparedTo(b.market_cap_rank);
      }));
      return;
    }
    setArrow('↓');
    setRatiosData(ratiosData.sort((a, b) => {
      const r = new BigNumber(a.global_ratio);
      return r.comparedTo(b.global_ratio) * -1;
    }));
  };

  const list = ratiosData.map((r) => {
    if (r.global_ratio === '0' || r.global_ratio === '') {
      return;
    }
    const ratios = r.ratios.map((rr) => {
      let name = 'Global';
      if (rr.category === 'TOP_LONG_SHORT_ACCOUNT_RATIO') {
        name = 'Top Account';
      }
      if (rr.category === 'TOP_LONG_SHORT_POSITION_RATIO') {
        name = 'Top Position';
      }
      return (
        <div key={rr.category}>
          {name}: {rr.long_short_ratio} / {rr.long_account} / {rr.short_account} ; {Math.floor((Date.now() - rr.timestamp) / 1000 / 60)}M
        </div>
      );
    });

    return (
      <div key={r.symbol} className={style.item}>
        <div className={style.info}>
          <img src={r.image} alt={r.name} className={style.icon} />
          {r.name} / {r.contract}
        </div>
        <div>
          {r.current_price} / {r.high_24h} / {r.low_24h}
        </div>
        <div>
          {r.market_cap_rank} / {r.market_cap}
        </div>
        <div>
          {r.funding_rate}
        </div>
        {ratios}
      </div>
    );
  });

  return (
    <div className='container'>
      <main className='column main'>
        <button onClick={handleSort}>sort{arrow}</button>
        {list}
      </main>
      <aside className='column aside'>
        <Widget />
      </aside>
    </div>
  );
};

export default Index;
