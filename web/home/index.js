import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class Home extends Component {
  constructor(props) {
    super(props);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('home', 'layout');
  }

  componentDidMount() {
    const api = new API();
    api.home.index(function (resp) {
      console.info(resp);
    });
  }

  render() {
    return (
      <HomeView />
    );
  }
}

const HomeView = (match) => (
  <div>
    <div>Home</div>
    <Link to='/about'>About</Link>
  </div>
);

export default Home;
