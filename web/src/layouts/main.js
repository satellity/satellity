import style from './main.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Route, Link, Switch, Redirect } from 'react-router-dom';
import logoURL from '../assets/images/chat.svg';
import API from '../api/index.js'
import Home from '../home/index.js';
import UserEdit from '../users/edit.js';
import UserShow from '../users/show.js';
import TopicNew from '../topics/new.js';
import TopicShow from '../topics/show.js';

class MainLayout extends Component {
  constructor(props) {
    super(props);
    this.state = {p: encodeURIComponent(props.location.pathname)};
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('main', 'layout');
  }

  render() {
    return (
      <div className={style.container}>
        <Header />
        <div className='wrapper'>
          <Switch>
            <Route exact path='/' component={Home} />
            <Route exact path='/user/edit' component={UserEdit} />
            <Route path='/users/:id' component={UserShow} />
            <Route exact path='/topics/new' component={TopicNew} />
            <Route path='/topics/:id/edit' component={TopicNew} />
            <Route path='/topics/:id' component={TopicShow} />
            <Redirect to={`/404?p=${this.state.p}`} />
          </Switch>
        </div>
      </div>
    )
  }
}

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

export default MainLayout;
