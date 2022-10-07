import React, {useEffect, useState} from 'react';
import API from 'api/index.js';

import style from './index.module.scss';

const api = new API();

const Index = () => {
  const [sources, setSources] = useState([]);

  useEffect(() => {
    api.source.admin.index().then((resp) => {
      if (resp.error) {
        return;
      }
      setSources(resp.data);
    });
  }, []);

  const list = sources.map((s) => {
    return (
      <div key={s.source_id} className={style.source}>
        <a href={s.link}>{s.link}</a>
        <div className={style.meta}>
          {s.author} · {s.host} · {s.locality} · {s.wreck} · {s.updated_at}
        </div>
      </div>
    );
  });

  return (
    <>
      {list}
    </>
  );
};

export default Index;
