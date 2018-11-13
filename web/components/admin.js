import React from 'react';
import { Route, Link } from 'react-router-dom';

const AdminRoute = ({component: Component, ...rest}) => {
  return (
    <Route {...rest} render={matchProps => (
      <div>
        <header className='header navi'>
          <span className='brand'>GD</span>
        </header>
        <div className='container'>
          <Component {...matchProps} />
        </div>
      </div>
    )} />
  )
};

export default AdminRoute;
