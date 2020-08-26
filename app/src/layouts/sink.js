import style from './sink.module.scss';
import React from 'react';
import { Link } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import Config from '../components/config.js';

const NoMatch = ({ location }) => {
  const classes = document.body.classList.values();
  document.body.classList.remove(...classes);
  document.body.classList.add('not-found', 'layout');
  let params = new URLSearchParams(location.search);
  let p = params.get('p') || '/404';

  return (
    <div className={style.container}>
      <Helmet>
        <title>{`Page Not Found - ${Config.Name}`}</title>
        <meta name='description' content='Page Not Found' />
        <link rel="canonical" href={`${Config.Host}/404`} />
      </Helmet>
      <h3 className={style.body}>
         LOL! NO MATCH FOR <span className={style.path}>{p}</span>
        <div className={style.action}>
          <Link to='/'>Back to homepage</Link>
        </div>
      </h3>
    </div>
  )
};

export default NoMatch;
