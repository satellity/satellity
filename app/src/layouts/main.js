import style from './main.module.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Route, Link, Switch, Redirect } from 'react-router-dom';
import Config from '../components/config.js';
import API from '../api/index.js'
import Home from '../home/view.js';
import User from '../users/view.js';
import Topic from '../topics/view.js';
import Tool from '../tools/view.js';
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
            <Route exact path='/' component={Topic.Index} />
            <Route exact path='/avatar' component={Tool.Portrait} />
            <Route exact path='/dashboard' component={Home.Dashboard} />
            <Route exact path='/user/edit' component={User.Edit} />
            <Route path='/users/:id' component={User.Show} />
            <Route exact path='/topics/new' component={Topic.New} />
            <Route path='/topics/:id/edit' component={Topic.New} />
            <Route path='/topics/:id' component={Topic.Show} />
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
    let profile = <span className={style.navi} onClick={this.handleLoginClick}>Login</span>;
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
          <Link to='/' className={style.brand}>
            <FontAwesomeIcon icon={['fa', 'home']} />
          </Link>
          <div className={style.site}><span className={style.name}>{Config.Name}</span></div>
          {profile}
        </header>
          {this.state.logging && <Modal handleLoginClick={this.handleLoginClick} />}
      </div>
    )
  }
}

export default MainLayout;
