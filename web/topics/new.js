import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class TopicNew extends Component {
  constructor(props) {
    super(props);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('topic', 'layout');
  }

  render() {
    return (
      <NewView />
    )
  }
}

const NewView = () => (
  <div> Hello Topics</div>
);

export default TopicNew;
