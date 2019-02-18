import style from './main.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React from 'react';
import { Route, Link } from 'react-router-dom';
import logoURL from '../assets/images/chat.svg';
import API from '../api/index.js'

const MainRoute = ({component: Component, ...rest}) => {
  const classes = document.body.classList.values();
  document.body.classList.remove(...classes);
  document.body.classList.add('main', 'layout');

  return (
    <Route {...rest} render={matchProps => (
      <div className={style.container}>
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
      <Link to='/user/edit' className={style.navi}> {user.readMe().nickname} </Link>
    );
  }
  return (
    <header className={style.header}>
      <Link to='/' className={style.brand}>
        <img src={logoURL} className={style.logo} alt='GoDiscourse'/>
        <span className={style.pc}>GoDiscourse</span>
        <span className={style.mobile}>GD</span>
      </Link>
      <Link to='/topics/new' className={style.navi}>
        <FontAwesomeIcon icon={['fa', 'plus']} />
      </Link>
      {profile}
    </header>
  )
}

export default MainRoute;
