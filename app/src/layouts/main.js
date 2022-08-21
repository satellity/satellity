import style from './main.module.scss';
import React, {Component} from 'react';
import {Route, Routes} from 'react-router-dom';
import Header from './header.js';
import Home from '../home/view.js';
import User from '../users/view.js';
import Topic from '../topics/view.js';

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
          <Routes>
            <Route index path='/' element={<Topic.Index />} />
            <Route exact path='/dashboard' element={<Home.Dashboard />} />
            <Route exact path='/categories/:id' element={<Topic.Index />} />
            <Route exact path='/user/edit' element={<User.Edit />} />
            <Route path='/users/:id' element={<User.Show />} />
            <Route exact path='/topics/new' element={<Topic.New />} />
            <Route path='/topics/:id/edit' element={<Topic.New />} />
            <Route path='/topics/:id' element={<Topic.Show />} />
          </Routes>
        </div>
      </div>
    );
  }
}

export default MainLayout;
