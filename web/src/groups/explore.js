import style from './explore.scss';
import React, { Component } from 'react';
import API from '../api/index.js';

class Explore extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {groups: []};
  }

  componentDidMount() {
    this.api.group.index().then((data) => {
      this.setState({groups: data});
    });
  }

  render() {
    return (
      <div className='wrapper container'>
        Group List
      </div>
    )
  }
}

export default Explore;
