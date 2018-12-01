import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../../api/index.js';

class AdminCategoryNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.state = {name: '', position: '', description: ''};
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
    this.api.category.create(this.state, (resp) => {
      history.push('/admin/categories');
    })
  }

  render() {
    return (
      <CategoryNew onSubmit={this.handleSubmit} onChange={this.handleChange} state={this.state} />
    )
  }
}

const CategoryNew = ({onSubmit, onChange, state}) => (
  <div className='admin categories'>
    <h1 className='welcome'>
      Create a new category
    </h1>
    <div className='panel'>
      <form onSubmit={(e) => onSubmit(e)}>
        <div>
          <label name='name'>Name *</label>
          <input type='text' name='name' value={state.name} autoComplete='off' onChange={(e) => onChange(e)} />
        </div>
        <div>
          <label name='position'>Position</label>
          <input type='number' name='position' value={state.position} pattern="[0-9]{10}" autoComplete='off' onChange={(e) => onChange(e)} />
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
  </div>
)

export default AdminCategoryNew;
