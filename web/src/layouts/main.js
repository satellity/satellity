import style from './main.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Route, Link, Switch, Redirect } from 'react-router-dom';
import Config from '../components/config.js';
import API from '../api/index.js'
import User from '../users/view.js';
import Topic from '../topics/view.js';
import Group from '../groups/view.js';

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
            <Route exact path='/' component={Topic.Index} />
            <Route exact path='/user/edit' component={User.Edit} />
            <Route path='/users/:id' component={User.Show} />
            <Route exact path='/community' component={Topic.Index} />
            <Route exact path='/topics/new' component={Topic.New} />
            <Route path='/topics/:id/edit' component={Topic.New} />
            <Route path='/topics/:id' component={Topic.Show} />
            <Route exact path='/groups' component={Group.Explore} />
            <Route exact path='/groups/new' component={Group.New} />
            <Route path='/groups/:id' component={Group.Show} />
            <Route exact path='/groups/:id/edit' component={Group.New} />
            <Redirect to={`/404?p=${this.state.p}`} />
          </Switch>
        </div>
      </div>
    )
  }
}

const Header = () => {
  const user = new API().user;
  let profile;
  if (user.loggedIn()) {
    profile = (
      <Link to='/user/edit' className={`${style.navi} ${style.user}`}> {user.readMe().nickname} </Link>
    );
  }
  return (
    <header className={style.header}>
      <Link to='/' className={style.brand}>
        <FontAwesomeIcon icon={['fa', 'home']} />
      </Link>
      <div className={style.site}><span className={style.name}>{Config.Name}</span></div>
      <Link to='/groups' className={style.navi}>
        Groups
      </Link>
      <Link to='/community' className={style.navi}>
        Community
      </Link>
      <Link to='/topics/new' className={style.navi}>
        <FontAwesomeIcon icon={['fa', 'plus']} />
      </Link>
      {profile}
    </header>
  )
}

export default MainLayout;
