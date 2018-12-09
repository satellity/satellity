import style from './main.scss';
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
        <div className='wrapper'>
          <Component {...matchProps} />
        </div>
      </div>
    )} />
  )
};

const Header = () => {
  const user = new API().user;
  let profile = '';
  if (user.loggedIn()) {
    profile = (
      <Link to='/user/edit' className={style.navi}> {user.me().nickname} </Link>
    );
  }
  return (
    <header className={style.header}>
      <Link to='/' className={style.brand}>
        <img src={logoURL} className={style.logo} alt='GoDiscourse'/>
        GoDiscourse
      </Link>
      <Link to='/topics/new' className={style.navi}>
        <FontAwesomeIcon icon={['fa', 'plus']} />
      </Link>
      {profile}
    </header>
  )
}

export default MainRoute;
