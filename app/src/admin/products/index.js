import style from './index.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../../api/index.js';

class Index extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {products: []};
  }

  componentDidMount() {
    this.api.product.index('').then((resp) => {
      if (resp.error) {
        return
      }
      this.setState({products: resp.data});
    });
  }

  render() {
    const state = this.state;

    const listProducts = state.products.map((product) => {
      return (
        <li key={product.product_id}>
          <Link to={`/products/${product.short_id}`}>{product.name}</Link>
          <div>
            {product.user.nickname} | {product.tags.join(', ')} | <a href={product.source} target='_blank' rel='noopener noreferrer'> External </a>
          </div>
          <div className={style.time}>
            {product.product_id} | {product.created_at} | <Link to={`/admin/products/${product.short_id}/edit`} >EIDT</Link>
          </div>
        </li>
      )
    });

    return (
      <div>
        <h1 className='welcome'>
            Here is the list of products.  <Link to='/admin/products/new' className='new'>Create New Product</Link>
        </h1>
        <div className='panel'>
          <ul>
              {listProducts}
          </ul>
        </div>
      </div>
    );
  }
}

export default Index;
