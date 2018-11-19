import './header.scss';
import React from 'react';
import { Route, Link } from 'react-router-dom';

const MainRoute = ({component: Component, ...rest}) => {
  return (
    <Route {...rest} render={matchProps => (
      <div>
        <div className='container'>
          <Component {...matchProps} />
        </div>
      </div>
    )} />
  )
};

export default MainRoute;
