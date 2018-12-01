import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class TopicNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {title: '', category_id: '', body: '', categories: []};
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

  componentDidMount() {
    const self = this;
    this.api.category.index((resp) => {
      let category_id = '';
      if (resp.data.length > 0) {
        category_id = resp.data[0].category_id;
      }
      self.setState({categories: resp.data, category_id: category_id});
    });
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
    const history = this.props.history;
    const data = {title: this.state.title, body: this.state.body, category_id: this.state.category_id};
    this.api.topic.create(data, (resp) => {
      history.push('/');
    });
  }

  render() {
    return (
      <View onSubmit={this.handleSubmit} onChange={this.handleChange} state={this.state} />
    )
  }
}

// TODO jsx editor format
const View = ({onSubmit, onChange, state}) => {
  const categories = state.categories.map((c) => {
    return (
      <option value={c.category_id} key={c.category_id}>{c.name}</option>
    )
  });

  return (
    <div>
      <h2>Create a new topic</h2>
      <form onSubmit={(e) => onSubmit(e)}>
        <div>
          <label name='name'>Title *</label>
          <input type='text' name='title' value={state.title} autoComplete='off' onChange={(e) => onChange(e)} />
        </div>
        <div>
          <label name='name'>Category</label>
          <div className='select'>
            <select name='category_id' value={state.category_id} onChange={(e) => onChange(e)}>
              {categories}
            </select>
          </div>
        </div>
        <div>
          <label name='body'>Body *</label>
          <textarea type='text' name='body' value={state.body} onChange={(e) => onChange(e)} />
        </div>
        <div className='action'>
          <input type='submit' value='SUBMIT' />
        </div>
      </form>
    </div>
  )
};

export default TopicNew;
