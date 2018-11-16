import './header.scss';
import React from 'react';
import { Route, Link } from 'react-router-dom';

const MainRoute = ({component: Component, ...rest}) => {
  return (
    <Route {...rest} render={matchProps => (
      <div>
        <header className='header navi'>
          <span className='brand'>GD</span>
          <div className='actions'>
            <Link to='/sign_in'>Sign In</Link>
          </div>
        </header>
        <div className='container'>
          <Component {...matchProps} />
        </div>
      </div>
    )} />
  )
};

export default MainRoute;
