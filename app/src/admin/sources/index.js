import React, {useEffect, useState} from 'react';
import {Link} from 'react-router-dom';
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

  const handleDelete = (e, id, title) => {
    e.preventDefault();
    const c = window.confirm(`Delete: ${title}`);
    if (c) {
      api.source.admin.delete(id).then((resp) => {
        if (resp.error) {
          return;
        }

        const ss = sources.filter((s) => {
          return s.source_id !== id;
        });
        setSources(ss);
      });
    }
  };

  const list = sources.map((s) => {
    return (
      <div key={s.source_id} className={style.source}>
        <a href={s.link}>{s.link}</a>
        &nbsp; &nbsp;
        <Link to='' onClick={(e) => handleDelete(e, s.source_id, s.link)} >❌</Link>
        <div className={style.meta}>
          {s.author} · {s.host} · {s.locality} · {s.wreck > 0 ? `errors ${s.wreck} · `: ''}{s.updated_at}
        </div>
        <div className={style.meta}>
          last publish at: {s.publish_at}
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
