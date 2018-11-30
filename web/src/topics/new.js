import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class TopicNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {title: '', body: ''}
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    const classes = document.body.classList.values();
    document.body.classList.remove(...classes);
    document.body.classList.add('topic', 'layout');
    // TODO handle authentication
    if (!this.api.user.loggedIn()) {
      props.history.push('/');
    }
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
  }

  render() {
    return (
      <View onSubmit={this.handleSubmit} onChange={this.handleChange} state={this.state} />
    )
  }
}

// TODO jsx editor format
const View = ({onSubmit, onChange, state}) => (
  <div>
    <h2>Create a new topic</h2>
    <form onSubmit={(e) => onSubmit(e)}>
      <div>
        <label name='name'>Title *</label>
        <input type='text' name='title' value={state.title} autoComplete='off' onChange={(e) => onChange(e)} />
      </div>
      <div>
        <label name='description'>Description *</label>
        <textarea type='text' name='description' value={state.description} onChange={(e) => onChange(e)} />
      </div>
      <div className='action'>
        <input type='submit' value='SUBMIT' />
      </div>
    </form>
  </div>
);

export default TopicNew;
