import React, {useEffect, useState} from 'react';
import {Link} from 'react-router-dom';
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

  const handleDelete = (e, id, title) => {
    e.preventDefault();
    const c = window.confirm(`Delete: ${title}`);
    if (c) {
      api.gist.admin.delete(id).then((resp) => {
        if (resp.error) {
          return;
        }

        const gs = gists.filter((g) => {
          return g.gist_id !== id;
        });
        setGists(gs);
      });
    }
  };

  const list = gists.map((g) => {
    return (
      <div key={g.gist_id} className={style.gist}>
        <a href={g.link}>{g.title}</a>
        &nbsp; &nbsp; &nbsp; &nbsp;
        <Link to='' onClick={(e) => handleDelete(e, g.gist_id, g.title)} >DELETE</Link>
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
