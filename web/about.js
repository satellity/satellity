import React, { Component } from 'react';
import { Link } from "react-router-dom";

class About extends Component {
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
