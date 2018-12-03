import './admin.scss';
import React, { Component } from 'react';
import { Route, Link } from 'react-router-dom';
import constants from '../components/constants.js';
import Index from './index.js';
import Categories from './categories/index.js';
import CategoriesNew from './categories/new.js';
import CategoriesEdit from './categories/edit.js';
import API from '../api/index.js';

class AdminRoute extends Component {
  constructor(props) {
    super(props);
    this.state = {site: constants.site};
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('admin', 'layout');
    if (!new API().user.isAdmin()) {
      props.history.push('/');
    }
  }

  render() {
    const match = this.props.match;
    return (
      <div>
        <header className='header'>
          <Link to='/' className='brand'> &larr; Back to {this.state.site}</Link>
          <Link to='/admin' className='navi'>Dashboard</Link>
          <Link to='/admin/categories' className='navi'>Categories</Link>
        </header>
        <div className='bg-container'>
          <div className='wrapper'>
            <Route exact path={`${match.url}`} component={Index} />
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
