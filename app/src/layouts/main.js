import style from './main.module.scss';
import React, {Component} from 'react';
import {Routes, Outlet} from 'react-router-dom';
import Header from './header.js';

class MainLayout extends Component {
  constructor(props) {
    super(props);

    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    this.state = {p: ''};
  }

  render() {
    return (
      <div className={style.container}>
        <Header />
        <div className='wrapper'>
          <Outlet />
          <Routes>
          </Routes>
        </div>
      </div>
    );
  }
}

export default MainLayout;
