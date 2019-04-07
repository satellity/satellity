import './admin.scss';
import React, { Component } from 'react';
import { Route, Link, Redirect } from 'react-router-dom';
import Config from '../components/config.js';
import Index from './index.js';
import Categories from './categories/index.js';
import Users from './users/index.js';
import Topics from './topics/index.js';
import CategoriesNew from './categories/new.js';
import CategoriesEdit from './categories/edit.js';
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
        <Redirect to={{ pathname: "/" }} />
      )
    }

    const match = this.props.match;
    return (
      <div>
        <header className='header'>
          <Link to='/' className='brand'> &larr; <span className='only-pc'>Back to {this.state.site}</span></Link>
          <Link to='/admin' className='navi'>Dashboard</Link>
          <Link to='/admin/users' className='navi'>Users</Link>
          <Link to='/admin/topics' className='navi'>Topics</Link>
          <Link to='/admin/categories' className='navi'>Categories</Link>
        </header>
        <div className='bg-container'>
          <div className='wrapper'>
            <Route exact path={`${match.url}`} component={Index} />
            <Route exact path={`${match.url}/users`} component={Users} />
            <Route exact path={`${match.url}/topics`} component={Topics} />
            <Route exact path={`${match.url}/categories`} component={Categories} />
            <Route exact path={`${match.url}/categories/new`} component={CategoriesNew} />
            <Route path={`${match.url}/categories/:id/edit`} component={CategoriesEdit} />
          </div>
        </div>
      </div>
    )
  }
}

export default AdminRoute;
