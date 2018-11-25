import './header.scss';
import React from 'react';
import { Route, Link } from 'react-router-dom';

const MainRoute = ({component: Component, ...rest}) => {
  return (
    <Route {...rest} render={matchProps => (
      <div>
        <Header />
        <div className='container'>
          <Component {...matchProps} />
        </div>
      </div>
    )} />
  )
};

const Header = () => {
  return (
    <header className='app header'>
      <Link to='/' className='brand'>
        SUNTIN
      </Link>
    </header>
  )
}

export default MainRoute;
