import React, { Component } from 'react';
import { Route, Link } from 'react-router-dom';
import About from '../about.js';
import API from '../api/index.js';

class AdminRoute extends Component {
  constructor(props) {
    super(props);
    if (!new API().user.isAdmin()) {
      props.history.push('/');
    }
  }

  render() {
    const match = this.props.match;
    return (
      <div>
        <header className='header navi'>
          <span className='brand'>GD Admin</span>
        </header>
        <div className='container'>
          <Route path={`${match.url}/about`} component={About} />
        </div>
      </div>
    )
  }
}

export default AdminRoute;
