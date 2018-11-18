import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../../api/index.js';

class AdminCategory extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
  }

  componentDidMount() {
    this.api.category.adminIndex(function(resp) {
      console.info(resp);
    });
  }

  render() {
    return (
      <div>
        <h1 className='welcome'>
          It is used to categorize topics.
          <Link to='admin/categories/new' className='new'>New Category</Link>
        </h1>
        <div className='panel'>
        </div>
      </div>
    )
  }
}

export default AdminCategory;
