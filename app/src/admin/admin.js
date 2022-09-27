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

  return (
    <>
      <header className={style.header}>
        <Link to='/' className={style.brand}> &larr; <span className='only-pc'>Back to {Config.Name}</span></Link>
        <Link to='/admin' className={style.navi}>Dashboard</Link>
        <Link to='/admin/users' className={style.navi}>Users</Link>
        <Link to='/admin/topics' className={style.navi}>Topics</Link>
        <Link to='/admin/comments' className={style.navi}>Comments</Link>
        <Link to='/admin/categories' className={style.navi}>Categories</Link>
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
