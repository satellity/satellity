import React, {useEffect, useState} from 'react';
import API from 'api/index.js';

import style from './index.module.scss';

const api = new API();

const Index = () => {
  const [gists, setGists] = useState([]);

  useEffect(() => {
    api.gist.admin.index().then((resp) => {
      if (resp.error) {
        return;
      }
      setGists(resp.data);
    });
  }, []);

  const list = gists.map((g) => {
    return (
      <div key={g.gist_id} className={style.gist}>
        {g.title}
        <div className={style.meta}>
          {g.source.author} · {g.source.host} · {g.publish_at}
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
