import './admin.scss';
import React, {Component} from 'react';
import {Route, Link, Navigate} from 'react-router-dom';
import Config from '../components/config.js';
import Index from './index.js';
import Category from './categories/view.js';
import Users from './users/index.js';
import Topics from './topics/index.js';
import Comments from './comments/index.js';
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
    if (!this.api.user.isAdmin()) {
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
            <Route index element={Index} />
            <Route exact path={`/users`} element={<Users />} />
            <Route exact path={`/topics`} element={<Topics />} />
            <Route exact path={`/comments`} element={<Comments />} />
            <Route exact path={`/categories`} element={<Category.Index />} />
            <Route exact path={`/categories/new`} element={<Category.New />} />
            <Route path={`/categories/:id/edit`} element={<Category.Edit />} />
          </div>
        </div>
      </div>
    );
  }
}

export default AdminRoute;
