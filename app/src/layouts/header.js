import React, {useState} from 'react';
import {Link} from 'react-router-dom';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import API from 'api/index.js';
import Config from 'components/config.js';
import Login from './login.js';

import logo from 'assets/images/logo.svg';
import style from './main.module.scss';

const Header = () => {
  const [logging, setLogging] = useState(false);

  const handleLoginClick = (e) => {
    const n = e.target.className;
    if (!(n.includes('close') || n.includes('modal') || n.includes('navi'))) {
      return;
    }
    setLogging(!logging);
  };

  const user = new API().user;
  let profile = <span className={style.navi} onClick={handleLoginClick}>Login</span>;
  if (user.loggedIn()) {
    profile = (
      <div className={style.navis}>
        <Link to='/topics/new' className={`${style.navi}`}> <FontAwesomeIcon icon={['fa', 'plus']} /> </Link>
        <Link to='/user/edit' className={`${style.navi} ${style.user}`}> {user.local().nickname} </Link>
      </div>
    );
  }

  return (
    <div>
      <header className={style.header}>
        <Link className={style.site} to='/'>
          <img className={style.logo} src={logo} alt={Config.Name} />
          <span className={style.name}>{Config.Name}</span>
        </Link>

        <div className={style.menus}>
          <Link className={`${style.menu} ${window.location.pathname === '/' ? style.current : ''}` } to='/'>
            Home
          </Link>
        </div>
        {profile}
      </header>
      {logging && <Login handleLoginClick={handleLoginClick} />}
    </div>
  );
};

export default Header;
