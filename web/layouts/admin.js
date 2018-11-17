import '../admin/admin.scss';
import React, { Component } from 'react';
import { Route, Link } from 'react-router-dom';
import Index from '../admin/index.js'
import About from '../about.js';
import API from '../api/index.js';

class AdminRoute extends Component {
  constructor(props) {
    super(props);
    if (!new API().user.isAdmin()) {
      props.history.push('/');
    }
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('admin', 'layout');
  }

  render() {
    const match = this.props.match;
    return (
      <div>
        <header className='header navi'>
          <Link to='/' className='brand'>FunYeah</Link>
          <Link to='/admin'>Dashboard</Link>
        </header>
        <div className='container'>
          <Route exact path={`${match.url}`} component={Index} />
          <Route path={`${match.url}/about`} component={About} />
        </div>
      </div>
    )
  }
}

export default AdminRoute;
