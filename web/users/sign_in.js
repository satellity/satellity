import './sign_in.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';

class SignIn extends Component {
  constructor(props) {
    super(props);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('sign_in', 'layout');
  }

  render() {
    // TODO replace with production githubClientId
    let githubClientId = '';
    if (process.env.NODE_ENV === 'development') {
      githubClientId = '03e10a9b62b4533e65b5';
    }
    return (
      <SignInView client_id={githubClientId} />
    );
  }
}

const SignInView = (props, match) => (
  <div>
    <h1 className='brand'><Link to='/'>Suntin</Link></h1>
    <div>
      <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${props.client_id}`} className='btn primary'>Sign in with GitHub</a>
    </div>
  </div>
);

export default SignIn;
