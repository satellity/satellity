import style from './main.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Route, Link, Switch, Redirect } from 'react-router-dom';
import Config from '../components/config.js';
import API from '../api/index.js'
import Home from '../home/view.js';
import User from '../users/view.js';
import Topic from '../topics/view.js';
import Group from '../groups/view.js';
import Message from '../messages/view.js';
import Modal from './modal.js';

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
            <Route exact path='/' component={Home.Index} />
            <Route exact path='/dashboard' component={Home.Dashboard} />
            <Route exact path='/user/edit' component={User.Edit} />
            <Route exact path='/user/groups' component={Group.List} />
            <Route path='/users/:id' component={User.Show} />
            <Route exact path='/community' component={Topic.Index} />
            <Route exact path='/topics/new' component={Topic.New} />
            <Route path='/topics/:id/edit' component={Topic.New} />
            <Route path='/topics/:id' component={Topic.Show} />
            <Route exact path='/groups' component={Group.Explore} />
            <Route exact path='/groups/new' component={Group.New} />
            <Route exact path='/groups/:id/edit' component={Group.New} />
            <Route exact path='/groups/:id/members' component={Group.Members} />
            <Route exact path='/groups/:id' component={Group.Show} />
            <Route exact path='/groups/:id/messages' component={Message.Index} />
            <Redirect to={`/404?p=${this.state.p}`} />
          </Switch>
        </div>
      </div>
    )
  }
}

class Header extends Component {
  constructor(props) {
    super(props);

    this.state = {logging: false};
    this.handleLoginClick = this.handleLoginClick.bind(this);
  }

  handleLoginClick(e) {
    this.setState({logging: !this.state.logging});
  }

  render() {
    const user = new API().user;
    let profile = <a className={style.navi} onClick={this.handleLoginClick}>Login</a>;
    if (user.loggedIn()) {
      profile = (
        <Link to='/user/edit' className={`${style.navi} ${style.user}`}> {user.local().nickname} </Link>
      );
    }

    return (
      <div>
        <header className={style.header}>
          <Link to='/' className={style.brand}>
            <FontAwesomeIcon icon={['fa', 'home']} />
          </Link>
          <div className={style.site}><span className={style.name}>{Config.Name}</span></div>
          <Link to='/groups' className={style.navi}>
              {i18n.t('group.name')}
          </Link>
          <Link to='/community' className={style.navi}>
              {i18n.t('community.name')}
          </Link>
            {profile}
        </header>
          {this.state.logging && <Modal handleLoginClick={this.handleLoginClick} />}
      </div>
    )
  }
}

export default MainLayout;
