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
    We are a team of two developers who love programming.
  </div>
);

export default About;
