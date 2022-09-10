import './admin.scss';
import React, {Component} from 'react';
import {Link, Navigate, Outlet} from 'react-router-dom';
import Config from '../components/config.js';
import API from '../api/index.js';

class AdminRoute extends Component {
  constructor(props) {
    super(props);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('admin', 'layout');
    this.api = new API();
    this.state = {site: Config.Name};
  }

  render() {
    if (!this.api.me.isAdmin()) {
      return (
        <Navigate to="/" replace />
      );
    }

    return (
      <div>
        <header className='header'>
          <Link to='/' className='brand'> &larr; <span className='only-pc'>Back to {this.state.site}</span></Link>
          <Link to='/admin' className='navi'>Dashboard</Link>
          <Link to='/admin/users' className='navi'>Users</Link>
          <Link to='/admin/topics' className='navi'>Topics</Link>
          <Link to='/admin/comments' className='navi'>Comments</Link>
          <Link to='/admin/categories' className='navi'>Categories</Link>
        </header>
        <div className='bg-container'>
          <div className='wrapper'>
            <Outlet />
          </div>
        </div>
      </div>
    );
  }
}

export default AdminRoute;
