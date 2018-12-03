import './header.scss';
import style from './header.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React from 'react';
import { Route, Link } from 'react-router-dom';
import logoURL from '../assets/images/logo.png';
import API from '../api/index.js'

const MainRoute = ({component: Component, ...rest}) => {
  return (
    <Route {...rest} render={matchProps => (
      <div>
        <Header />
        <div className='wrapper section app container'>
          <Component {...matchProps} />
        </div>
      </div>
    )} />
  )
};

const Header = () => {
  const user = new API().user;
  let link = (<Link to='/sign_in' className={style.navi}> SignIn </Link>);
  if (user.loggedIn()) {
    link = (
      <Link to='/sign_in' className={style.navi}> {user.me().nickname} </Link>
    );
  }
  return (
    <header className={style.header}>
      <Link to='/' className={style.brand}>
        <img src={logoURL} className={style.logo} alt='SUNTIN'/>
        SUNTIN
      </Link>
      <Link to='/topics/new' className={style.navi}>
        <FontAwesomeIcon icon={['fa', 'plus']} />
      </Link>
      {link}
    </header>
  )
}

export default MainRoute;
