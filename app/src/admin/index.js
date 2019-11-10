import React, { Component } from 'react';
import Config from '../components/config.js';

class Index extends Component {
  render() {
    return (
      <h1 className='welcome'>
          This is the Dashboard for {Config.Name}.
      </h1>
    )
  }
}

export default Index;
