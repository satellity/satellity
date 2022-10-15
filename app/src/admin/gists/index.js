import React, {useEffect, useState} from 'react';
import {useParams, Link} from 'react-router-dom';
import API from 'api/index.js';

import style from './index.module.scss';

const api = new API();

const Index = () => {
  const {genre} = useParams();

  const [gists, setGists] = useState([]);

  useEffect(() => {
    let request = api.gist.admin.index();
    if (!!genre) {
      request = api.gist.admin.genres(genre);
    }
    request.then((resp) => {
      if (resp.error) {
        return;
      }
      setGists(resp.data);
    });
  }, [genre]);


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
      <div>
        <Link to="/admin/gists" className={`${style.node} ${!genre ? style.current : ''}`}> News </Link>
        <Link to="/admin/genres/release" className={`${style.node} ${genre === 'release' ? style.current : ''}`}> Releases </Link>
      </div>
      {list}
    </>
  );
};

export default Index;
