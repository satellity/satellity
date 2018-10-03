import './sign_in.scss';
import React, { Component } from 'react';
import { Link } from "react-router-dom";

class SignIn extends Component {
  constructor(props) {
    super(props);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('sign_in', 'layout');
  }

  render() {
    return (
      <SignInView />
    );
  }
}

const SignInView = (match) => (
  <div>
    <h1 className="brand"><Link to="/">GD</Link></h1>
    <div className="slogan">
      A discourse like forum.
    </div>
    <div>
      <a href="javascript:;" className="button-warning pure-button button">Sign in with GitHub</a>
    </div>
  </div>
);

export default SignIn;
