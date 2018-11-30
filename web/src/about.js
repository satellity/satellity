import React, { Component } from 'react';
import { Link } from "react-router-dom";

class About extends Component {
  constructor(props) {
    super(props);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('about', 'layout');
  }

  render() {
    return (
      <AboutView />
    );
  }
}

const AboutView = () => (
  <div>
    <div>About</div>
    <Link to='/'>Home</Link>
  </div>
);

export default About;
