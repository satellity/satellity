import React from 'react';
import {Link, Navigate, Outlet} from 'react-router-dom';
import Config from 'components/config.js';
import API from '../api/index.js';

import style from './admin.module.scss';

const api = new API();

const Layout = () => {
  if (!api.me.isAdmin()) {
    return (
      <Navigate to="/" replace />
    );
  }

  const navis = [
    ['/admin', 'Dashboard'],
    ['/admin/users', 'Users'],
    ['/admin/topics', 'Topics'],
    ['/admin/comments', 'Comments'],
    ['/admin/categories', 'Categories'],
    ['/admin/gists', 'Gists'],
  ];

  const views = navis.map((n) => {
    return (
      <Link key={n[1]} to={n[0]} className={style.navi}>{n[1]}</Link>
    );
  });

  return (
    <>
      <header className={style.header}>
        <Link to='/' className={style.brand}> &larr; <span className='only-pc'>Back to {Config.Name}</span></Link>
        {views}
      </header>
      <div className={style.container}>
        <div className={style.wrapper}>
          <Outlet />
        </div>
      </div>
    </>
  );
};

export default Layout;
