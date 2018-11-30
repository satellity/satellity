import './header.scss';
import React from 'react';
import { Route, Link } from 'react-router-dom';
import logoURL from '../assets/images/logo.png';
import API from '../api/index.js'

const MainRoute = ({component: Component, ...rest}) => {
  return (
    <Route {...rest} render={matchProps => (
      <div>
        <Header />
        <div className='app container'>
          <Component {...matchProps} />
        </div>
      </div>
    )} />
  )
};

const Header = () => {
  const user = new API().user;
  let link = (<Link to='/sign_in' className='navi'> SignIn </Link>);
  if (user.loggedIn()) {
    link = (<Link to='/sign_in' className='navi'> {user.me().nickname} </Link>);
  }
  return (
    <header className='app header'>
      <Link to='/' className='brand'>
        <img src={logoURL} className='logo' alt=''/>
        SUNTIN
      </Link>
      {link}
    </header>
  )
}

export default MainRoute;
