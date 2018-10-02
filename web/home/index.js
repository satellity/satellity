import React, { Component } from 'react';
import { Link } from "react-router-dom";

class Home extends Component {
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
